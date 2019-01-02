package di

import (
	"context"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"

	"gitlab.com/proemergotech/dliver-project-skeleton/app/client/centrifugo/proto/apiproto"
	"google.golang.org/grpc"

	"github.com/json-iterator/go"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/client/redis"

	"github.com/go-playground/validator"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	jconfig "github.com/uber/jaeger-client-go/config"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/event"
	"gitlab.com/proemergotech/geb-client-go/geb"
	"gitlab.com/proemergotech/geb-client-go/geb/rabbitmq"
	"gitlab.com/proemergotech/log-go/geblog"
	"gitlab.com/proemergotech/trace-go"
	"gitlab.com/proemergotech/trace-go/gebtrace"

	"gitlab.com/proemergotech/dliver-project-skeleton/app"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/action"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/apierr"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/config"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/rest"
	"gitlab.com/proemergotech/log-go"
	"gitlab.com/proemergotech/log-go/echolog"
	"gitlab.com/proemergotech/log-go/jaegerlog"
)

type Container struct {
	RestServer       *rest.Server
	EventServer      *event.Server
	redisClient      *redis.Client
	centrifugeClient apiproto.CentrifugeClient
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

	// Init Tracer
	tracer, closer, err := newTracer(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize Jaeger Tracer")
	}
	c.traceCloser = closer
	opentracing.SetGlobalTracer(tracer)

	// Init GEB queue
	gebQueue, err := newGebQueue(cfg, tracer)
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize geb queue")
	}
	c.gebCloser = gebQueue

	// Init REDIS client
	c.redisClient, err = newRedis(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize redis client")
	}

	// Init Centrifuge Client
	c.centrifugeClient, err = newCentrifugeClient(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize centrifuge client")
	}

	// Init Validator
	validate := newValidator()

	// Init ECHO rest server engine
	echoEngine := newEcho(validate)

	// Init service core
	core := app.NewCore(c.centrifugeClient)

	// Init service actions
	acts := action.NewActions(core, validate)

	// Init REST server
	c.RestServer = rest.NewServer(
		cfg.Port,
		echoEngine,
		rest.NewRouter(
			echoEngine,
			acts,
			tracer,
		),
	)

	// Init EVENT server
	c.EventServer = event.NewServer(event.NewRouter(
		gebQueue,
		validate,
	))

	return c, nil
}

func newTracer(cfg *config.Config) (opentracing.Tracer, io.Closer, error) {
	// init Jaeger tracer
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

func newRedis(cfg *config.Config) (*redis.Client, error) {
	redisPool, err := redis.NewRedisPool(cfg)
	if err != nil {
		return nil, err
	}

	redisJSON := jsoniter.Config{
		SortMapKeys:            true,
		ValidateJsonRawMessage: true,
		OnlyTaggedField:        true,
		TagKey:                 "redis",
	}.Froze()

	return redis.NewClient(redisPool, redisJSON), nil
}

func newCentrifugeClient(cfg *config.Config) (apiproto.CentrifugeClient, error) {
	grpcConn, err := grpc.Dial(fmt.Sprintf("%v:%v", cfg.CentrifugoHost, cfg.CentrifugoGrpcPort), grpc.WithInsecure())
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize centrifuge connection")
	}
	return apiproto.NewCentrifugeClient(grpcConn), nil
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

func (c *Container) Close() {
	err := c.gebCloser.Close()
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
