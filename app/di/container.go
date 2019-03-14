package di

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	jaeger "github.com/uber/jaeger-client-go"
	jconfig "github.com/uber/jaeger-client-go/config"
	centrifuge "gitlab.com/proemergotech/centrifuge-client-go"
	"gitlab.com/proemergotech/centrifuge-client-go/api"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/apierr"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/config"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/event"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/rest"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/storage"
	"gitlab.com/proemergotech/geb-client-go/geb"
	"gitlab.com/proemergotech/geb-client-go/geb/rabbitmq"
	log "gitlab.com/proemergotech/log-go"
	"gitlab.com/proemergotech/log-go/echolog"
	"gitlab.com/proemergotech/log-go/geblog"
	"gitlab.com/proemergotech/log-go/gentlemanlog"
	"gitlab.com/proemergotech/log-go/jaegerlog"
	trace "gitlab.com/proemergotech/trace-go"
	"gitlab.com/proemergotech/trace-go/gebtrace"
	"gitlab.com/proemergotech/trace-go/gentlemantrace"
	gentleman "gopkg.in/h2non/gentleman.v2"
)

type Container struct {
	RestServer       *rest.Server
	EventServer      *event.Server
	redisClient      *storage.Redis
	centrifugeClient api.CentrifugeClient
	traceCloser      io.Closer
	gebCloser        io.Closer
}

type EchoValidator struct {
	validator *validator.Validate
}

func (cv *EchoValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)
	if err != nil {
		return apierr.Validation(err)
	}

	return nil
}

func NewContainer(cfg *config.Config) (*Container, error) {
	c := &Container{}

	tracer, closer, err := newTracer(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize Jaeger Tracer")
	}
	c.traceCloser = closer
	opentracing.SetGlobalTracer(tracer)

	gebQueue, err := newGebQueue(cfg, tracer)
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize geb queue")
	}
	c.gebCloser = gebQueue

	c.redisClient, err = newRedis(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize redis client")
	}

	c.centrifugeClient, err = centrifuge.New(cfg.CentrifugoHost, cfg.CentrifugoGrpcPort, 5*time.Second)
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize centrifuge client")
	}

	validate := newValidator()

	echoEngine := newEcho(validate)

	svc := service.NewService(
		c.centrifugeClient,
	)

	c.RestServer = rest.NewServer(
		cfg.Port,
		echoEngine,
		rest.NewController(
			echoEngine,
			svc,
			tracer,
		),
	)

	c.EventServer = event.NewServer(event.NewController(
		gebQueue,
		validate,
		svc,
	))

	return c, nil
}

func newTracer(cfg *config.Config) (opentracing.Tracer, io.Closer, error) {
	transport, err := jaeger.NewUDPTransport(
		fmt.Sprintf("%v:%v", cfg.TracerReporterLocalAgentHost, cfg.TracerReporterLocalAgentPort),
		8000,
	)
	if err != nil {
		return nil, nil, errors.Wrap(err, "couldn't create udp transport for jaeger")
	}

	tracerSamplerParam, err := strconv.ParseFloat(cfg.TracerSamplerParam, 64)
	if err != nil {
		return nil, nil, errors.Wrap(err, "couldn't load configuration for tracing")
	}

	trcConf := &jconfig.Configuration{
		Sampler: &jconfig.SamplerConfig{
			Type:  cfg.TracerSamplerType,
			Param: tracerSamplerParam,
			SamplingServerURL: fmt.Sprintf(
				"%v://%v:%v",
				cfg.TracerSamplerSamplingServerScheme,
				cfg.TracerSamplerSamplingServerHost,
				cfg.TracerSamplerSamplingServerPort,
			),
		},
		ServiceName: config.AppName,
	}
	return trcConf.NewTracer(
		jconfig.Logger(jaegerlog.NewJaegerLogger(log.GlobalLogger())),
		jconfig.Reporter(jaeger.NewRemoteReporter(transport, jaeger.ReporterOptions.Logger(jaegerlog.NewJaegerLogger(log.GlobalLogger())))),
	)
}

func newGebQueue(cfg *config.Config, tracer opentracing.Tracer) (*geb.Queue, error) {
	q := geb.NewQueue(
		rabbitmq.NewHandler(
			config.AppName,
			cfg.GebUsername,
			cfg.GebPassword,
			cfg.GebHost,
			cfg.GebPort,
			rabbitmq.Timeout(2*time.Second),
		),
		geb.JSONCodec(geb.UseTag("geb")),
	)

	opt, err := gebtrace.Trace(trace.Start)
	if err != nil {
		return nil, err
	}

	q.UsePublish(geblog.PublishMiddleware(log.GlobalLogger(), true))
	q.UsePublish(gebtrace.PublishMiddleware(tracer, log.GlobalLogger()))
	q.UseOnEvent(geblog.OnEventMiddleware(log.GlobalLogger(), true))
	q.UseOnEvent(gebtrace.OnEventMiddleware(tracer, log.GlobalLogger(), gebtrace.GenerateCorrelation(trace.NewCorrelation), opt))
	q.UseOnEvent(func(e *geb.Event, next func(*geb.Event) error) error {
		err := next(e)
		if err != nil {
			statusCode := apierr.StatusCode(err)
			if statusCode >= 400 && statusCode < 500 {
				log.Warn(e.Context(), err.Error(), "error", err)
			} else {
				log.Error(e.Context(), err.Error(), "error", err)
			}
		}

		return nil
	})
	q.OnError(func(err error) {
		err = errors.Wrap(err, "Geb connection error")
		log.Error(context.Background(), err.Error(), "error", err)

		go func() {
			time.Sleep(2 * time.Second)
			q.Reconnect()
		}()
	})

	return q, nil
}

func newRedis(cfg *config.Config) (*storage.Redis, error) {
	redisPool, err := storage.NewRedisPool(cfg)
	if err != nil {
		return nil, err
	}

	redisJSON := jsoniter.Config{
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		OnlyTaggedField:        true,
		TagKey:                 "redis",
	}.Froze()

	return storage.NewRedis(redisPool, redisJSON), nil
}

func newValidator() *validator.Validate {
	v := validator.New()

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			name = ""
		}

		return name
	})

	return v
}

func newEcho(validate *validator.Validate) *echo.Echo {
	e := echo.New()

	e.Use(echolog.DebugMiddleware(log.GlobalLogger(), true, true))
	e.Use(middleware.Recover())
	e.HTTPErrorHandler = rest.DLiveRHTTPErrorHandler
	e.Validator = &EchoValidator{validator: validate}

	return e
}

func newGentleman(scheme string, host string, port int, tracer opentracing.Tracer) *gentleman.Client {
	return gentleman.New().BaseURL(fmt.Sprintf("%v://%v:%v", scheme, host, port)).
		Use(gentlemantrace.Middleware(tracer, log.GlobalLogger())).
		Use(gentlemanlog.Middleware(log.GlobalLogger(), true, true))
}

func (c *Container) Close() {
	err := c.redisClient.Close()
	if err != nil {
		err = errors.Wrap(err, "Redis graceful close failed")
		log.Warn(context.Background(), err.Error(), "error", err)
	}

	err = c.gebCloser.Close()
	if err != nil {
		err = errors.Wrap(err, "GebQueue graceful close failed")
		log.Warn(context.Background(), err.Error(), "error", err)
	}

	err = c.traceCloser.Close()
	if err != nil {
		err = errors.Wrap(err, "Tracer graceful close failed")
		log.Warn(context.Background(), err.Error(), "error", err)
	}
}
