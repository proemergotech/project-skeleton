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
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo"
	"github.com/olivere/elastic"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	jaeger "github.com/uber/jaeger-client-go"
	jconfig "github.com/uber/jaeger-client-go/config"
	centrifuge "gitlab.com/proemergotech/centrifuge-client-go"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/config"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/event"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/rest"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/storage"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/validation"
	"gitlab.com/proemergotech/geb-client-go/geb"
	"gitlab.com/proemergotech/geb-client-go/geb/rabbitmq"
	log "gitlab.com/proemergotech/log-go"
	"gitlab.com/proemergotech/log-go/echolog"
	"gitlab.com/proemergotech/log-go/elasticlog"
	"gitlab.com/proemergotech/log-go/geblog"
	"gitlab.com/proemergotech/log-go/gentlemanlog"
	"gitlab.com/proemergotech/log-go/httplog"
	"gitlab.com/proemergotech/log-go/jaegerlog"
	"gitlab.com/proemergotech/trace-go/gebtrace"
	"gitlab.com/proemergotech/trace-go/gentlemantrace"
	yafuds "gitlab.com/proemergotech/yafuds-client-go/client"
	gentleman "gopkg.in/h2non/gentleman.v2"
)

type Container struct {
	RestServer    *rest.Server
	EventServer   *event.Server
	redisCloser   io.Closer
	traceCloser   io.Closer
	gebCloser     io.Closer
	yafudsCloser  io.Closer
	elasticClient *elastic.Client
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

	closer, err := newTracer(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize Jaeger Tracer")
	}
	c.traceCloser = closer

	gebQueue, err := newGebQueue(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize geb queue")
	}
	c.gebCloser = gebQueue

	redisStore, err := newRedisStore(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize redis client")
	}
	c.redisCloser = redisStore

	yafudsClient, err := newYafuds(cfg)
	if err != nil {
		return nil, err
	}
	c.yafudsCloser = yafudsClient

	validator := newValidator()

	echoEngine := newEcho(validator)

	svc := service.NewService(
		centrifugeClient,
		yafudsClient,
	)

	c.RestServer = rest.NewServer(
		newHTTPServer(
			echoEngine,
			cfg.Port,
		),
		rest.NewController(
			echoEngine,
			svc,
		),
	)

	c.EventServer = event.NewServer(
		event.NewController(
			gebQueue,
			validator,
			svc,
		),
	)

	return c, nil
}

func newTracer(cfg *config.Config) (io.Closer, error) {
	transport, err := jaeger.NewUDPTransport(
		fmt.Sprintf("%v:%v", cfg.TracerReporterLocalAgentHost, cfg.TracerReporterLocalAgentPort),
		8000,
	)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't create udp transport for jaeger")
	}

	tracerSamplerParam, err := strconv.ParseFloat(cfg.TracerSamplerParam, 64)
	if err != nil {
		return nil, errors.Wrap(err, "couldn't load configuration for tracing")
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

	tracer, closer, err := trcConf.NewTracer(
		jconfig.Logger(jaegerlog.NewJaegerLogger(log.GlobalLogger())),
		jconfig.Reporter(jaeger.NewRemoteReporter(transport, jaeger.ReporterOptions.Logger(jaegerlog.NewJaegerLogger(log.GlobalLogger())))),
	)
	if err != nil {
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)
	return closer, nil
}

func newElastic(cfg *config.Config) (*storage.Elastic, error) {
	httpClient := &http.Client{Transport: httplog.NewLoggingTransport(http.DefaultTransport, log.GlobalLogger(), "Elasticsearch: ", true, true)}
	elasticClient, err := elastic.NewClient(
		elastic.SetErrorLog(elasticlog.NewErrorLogger(log.GlobalLogger())),
		elastic.SetHealthcheck(false),
		elastic.SetHttpClient(httpClient),
		elastic.SetRetrier(elastic.NewBackoffRetrier(elastic.NewExponentialBackoff(100*time.Millisecond, 1*time.Second))),
		elastic.SetSniff(false),
		elastic.SetURL(fmt.Sprintf("%v://%v:%v", cfg.ElasticSearchScheme, cfg.ElasticSearchHost, cfg.ElasticSearchPort)),
	)
	if err != nil {
		return nil, errors.Wrap(err, "failed creating elastic client")
	}

	return storage.NewElastic(elasticClient), nil
}

func newGebQueue(cfg *config.Config) (*geb.Queue, error) {
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

	q.UsePublish(geb.RetryMiddleware())
	q.UsePublish(geblog.PublishMiddleware(log.GlobalLogger(), true))
	q.UsePublish(gebtrace.PublishMiddleware(opentracing.GlobalTracer(), log.GlobalLogger()))
	q.UseOnEvent(geb.RecoveryMiddleware())
	q.UseOnEvent(geblog.OnEventMiddleware(log.GlobalLogger(), true))
	q.UseOnEvent(gebtrace.OnEventMiddleware(opentracing.GlobalTracer(), log.GlobalLogger()))
	q.UseOnEvent(func(e *geb.Event, next func(*geb.Event) error) error {
		if err := next(e); err != nil {
			httpCode := schema.ErrorHTTPCode(err)
			if httpCode >= 400 && httpCode < 500 {
				log.Warn(e.Context(), err.Error(), "error", err)
			} else {
				log.Error(e.Context(), err.Error(), "error", err)
			}
		}

		return nil
	})
	if err := q.OnError(func(err error, reconnect func()) {
		err = errors.Wrap(err, "Geb connection error")
		log.Error(context.Background(), err.Error(), "error", err)

		go func() {
			time.Sleep(2 * time.Second)
			reconnect()
		}()
	}); err != nil {
		return nil, err
	}

	return q, nil
}

func newRedisStore(cfg *config.Config) (*storage.Redis, error) {
	redisPool, err := newRedisPool(
		cfg.RedisStorePoolIdleTimeout,
		cfg.RedisStorePoolMaxIdle,
		cfg.RedisStoreHost,
		cfg.RedisStorePort,
		cfg.RedisStoreDatabase,
	)
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

func newRedisPool(poolIdleTimeout string, poolMaxIdle int, host string, port int, database int) (*redis.Pool, error) {
	redisPoolIdleTimeout, err := time.ParseDuration(poolIdleTimeout)
	if err != nil {
		return nil, errors.Wrap(err, "invalid value for redis_pool_idle_timeout, must be duration")
	}

	return &redis.Pool{
		MaxIdle:     poolMaxIdle,
		IdleTimeout: redisPoolIdleTimeout,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", fmt.Sprintf("%v:%v", host, port), redis.DialDatabase(database))
		},
	}, nil
}

func newYafuds(cfg *config.Config) (yafuds.Client, error) {
	yafuds.SetTracer(opentracing.GlobalTracer())
	yafudsClient, err := yafuds.New(cfg.YafudsHost, cfg.YafudsPort, yafuds.Timeout(5*time.Second), yafuds.Retries(3))
	if err != nil {
		return nil, errors.Wrap(err, "unable to connect to Yafuds")
	}
	yafuds.SetLogger(log.GlobalLogger())

	return yafudsClient, nil
}

func newValidator() *validation.Validator {
	v := validator.New()

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		name := strings.SplitN(field.Tag.Get("json"), ",", 2)[0]

		if name == "-" {
			name = ""
		}

		return name
	})

	// TODO: remove example validation for enums:
	//err = v.RegisterValidation("enum_status", func(fl validator.FieldLevel) bool {
	//	return schema.Statuses[fl.Field().String()]
	//})
	//if err != nil {
	//	return nil, err
	//}

	return validation.NewValidator(v)
}

func newHTTPServer(echoEngine *echo.Echo, port int) *http.Server {
	return &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: echoEngine,
	}
}

func newEcho(validator *validation.Validator) *echo.Echo {
	e := echo.New()

	e.Use(echolog.DebugMiddleware(log.GlobalLogger(), true, true))
	e.Use(echolog.RecoveryMiddleware(log.GlobalLogger()))
	e.HTTPErrorHandler = rest.DLiveRHTTPErrorHandler
	e.Validator = validator

	return e
}

func newGentleman(scheme string, host string, port int) *gentleman.Client {
	return gentleman.New().BaseURL(fmt.Sprintf("%v://%v:%v", scheme, host, port)).
		Use(gentlemantrace.Middleware(opentracing.GlobalTracer(), log.GlobalLogger())).
		Use(gentlemanlog.Middleware(log.GlobalLogger(), true, true))
}

func (c *Container) Close() {
	c.elasticClient.Stop()

	if err := c.gebCloser.Close(); err != nil {
		err = errors.Wrap(err, "gebQueue graceful close failed")
		log.Warn(context.Background(), err.Error(), "error", err)
	}

	if err := c.redisCloser.Close(); err != nil {
		err = errors.Wrap(err, "redis graceful close failed")
		log.Warn(context.Background(), err.Error(), "error", err)
	}

	if err := c.traceCloser.Close(); err != nil {
		err = errors.Wrap(err, "tracer graceful close failed")
		log.Warn(context.Background(), err.Error(), "error", err)
	}

	if err := c.yafudsCloser.Close(); err != nil {
		err = errors.Wrap(err, "yafuds graceful close failed")
		log.Warn(context.Background(), err.Error(), "error", err)
	}
}
