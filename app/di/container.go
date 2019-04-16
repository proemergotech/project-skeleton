package di

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/gomodule/redigo/redis"
	"github.com/json-iterator/go"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/olivere/elastic"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"github.com/uber/jaeger-client-go"
	jconfig "github.com/uber/jaeger-client-go/config"
	"gitlab.com/proemergotech/centrifuge-client-go"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/config"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/event"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/rest"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/storage"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/validationerr"
	"gitlab.com/proemergotech/geb-client-go/geb"
	"gitlab.com/proemergotech/geb-client-go/geb/rabbitmq"
	"gitlab.com/proemergotech/log-go"
	"gitlab.com/proemergotech/log-go/echolog"
	"gitlab.com/proemergotech/log-go/elasticlog"
	"gitlab.com/proemergotech/log-go/geblog"
	"gitlab.com/proemergotech/log-go/gentlemanlog"
	"gitlab.com/proemergotech/log-go/httplog"
	"gitlab.com/proemergotech/log-go/jaegerlog"
	"gitlab.com/proemergotech/trace-go/gebtrace"
	"gitlab.com/proemergotech/trace-go/gentlemantrace"
	yclient "gitlab.com/proemergotech/yafuds-client-go/client"

	"gopkg.in/h2non/gentleman.v2"
)

type Container struct {
	RestServer    *rest.Server
	EventServer   *event.Server
	redisClient   *storage.Redis
	traceCloser   io.Closer
	gebCloser     io.Closer
	yafudsCloser  io.Closer
	elasticClient *elastic.Client
}

type EchoValidator struct {
	validator *validator.Validate
}

func (cv *EchoValidator) Validate(i interface{}) error {
	err := cv.validator.Struct(i)
	if err != nil {
		return validationerr.ValidationError{Err: err}.E()
	}

	return nil
}

func NewContainer(cfg *config.Config) (*Container, error) {
	c := &Container{}

	centrifuge.SetLogger(log.GlobalLogger())
	centrifugeClient, err := centrifuge.New(cfg.CentrifugoHost, cfg.CentrifugoGrpcPort, centrifuge.Timeout(5*time.Second))
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize centrifuge client")
	}

	e, err := newElastic(cfg)
	if err != nil {
		return nil, err
	}
	c.elasticClient = e.ElasticClient

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

	yafuds, err := newYafuds(cfg)
	if err != nil {
		return nil, err
	}
	c.yafudsCloser = yafuds

	validate := newValidator()

	echoEngine := newEcho(validate)

	svc := service.NewService(
		centrifugeClient,
		yafuds,
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

func newElastic(cfg *config.Config) (*storage.Elastic, error) {
	httpClient := &http.Client{Transport: httplog.NewLoggingTransport(http.DefaultTransport, log.GlobalLogger(), "Elasticsearch: ", true, true)}
	elasticClient, err := elastic.NewClient(
		elastic.SetErrorLog(elasticlog.NewErrorLogger(log.GlobalLogger())),
		elastic.SetHealthcheck(false),
		elastic.SetHttpClient(httpClient),
		elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewExponentialBackoff(100*time.Millisecond, 1*time.Second))),
		elastic.SetSniff(false),
		elastic.SetURL(cfg.ElasticAddress),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating elastic client")
	}

	return storage.NewElastic(elasticClient), nil
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

	q.UsePublish(geblog.PublishMiddleware(log.GlobalLogger(), true))
	q.UsePublish(gebtrace.PublishMiddleware(tracer, log.GlobalLogger()))
	q.UseOnEvent(geb.RecoveryMiddleware())
	q.UseOnEvent(geblog.OnEventMiddleware(log.GlobalLogger(), true))
	q.UseOnEvent(gebtrace.OnEventMiddleware(tracer, log.GlobalLogger()))
	q.UseOnEvent(func(e *geb.Event, next func(*geb.Event) error) error {
		err := next(e)
		if err != nil {
			httpCode := schema.ErrorHTTPCode(err)
			if httpCode >= 400 && httpCode < 500 {
				log.Warn(e.Context(), err.Error(), "error", err)
			} else {
				log.Error(e.Context(), err.Error(), "error", err)
			}
		}

		return nil
	})
	err := q.OnError(func(err error, reconnect func()) {
		err = errors.Wrap(err, "Geb connection error")
		log.Error(context.Background(), err.Error(), "error", err)

		go func() {
			time.Sleep(2 * time.Second)
			reconnect()
		}()
	})
	if err != nil {
		return nil, err
	}

	return q, nil
}

func newRedis(cfg *config.Config) (*storage.Redis, error) {
	redisPool, err := newRedisPool(cfg)
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

func newRedisPool(cfg *config.Config) (*redis.Pool, error) {
	redisPoolIdleTimeout, err := time.ParseDuration(cfg.RedisStorePoolIdleTimeout)
	if err != nil {
		return nil, errors.New("invalid value for redis_pool_idle_timeout, must be duration")
	}

	return &redis.Pool{
		MaxIdle:     cfg.RedisStorePoolMaxIdle,
		IdleTimeout: redisPoolIdleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%v:%v", cfg.RedisStoreHost, cfg.RedisStorePort), redis.DialDatabase(cfg.RedisStoreDatabase))
		},
	}, nil
}

func newYafuds(cfg *config.Config) (*yclient.Client, error) {
	yafudsClient, err := yclient.New(cfg.YafudsHost, cfg.YafudsPort, yclient.Timeout(5*time.Second), yclient.Retries(3))
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to Yafuds")
	}
	yclient.SetLogger(log.GlobalLogger())

	return yafudsClient, nil
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
	c.elasticClient.Stop()

	err := c.gebCloser.Close()
	if err != nil {
		err = errors.Wrap(err, "gebQueue graceful close failed")
		log.Warn(context.Background(), err.Error(), "error", err)
	}

	err = c.redisClient.Close()
	if err != nil {
		err = errors.Wrap(err, "redis graceful close failed")
		log.Warn(context.Background(), err.Error(), "error", err)
	}

	err = c.traceCloser.Close()
	if err != nil {
		err = errors.Wrap(err, "tracer graceful close failed")
		log.Warn(context.Background(), err.Error(), "error", err)
	}

	err = c.yafudsCloser.Close()
	if err != nil {
		err = errors.Wrap(err, "yafuds graceful close failed")
		log.Warn(context.Background(), err.Error(), "error", err)
	}
}
