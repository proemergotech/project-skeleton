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

	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/non-standard/validators"
	"github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
	"github.com/olivere/elastic"
	"github.com/opentracing/opentracing-go"
	"github.com/proemergotech/errors"
	"github.com/proemergotech/log/v3"
	"github.com/proemergotech/log/v3/echolog"
	"github.com/proemergotech/log/v3/elasticlog"
	"github.com/proemergotech/log/v3/httplog"
	"github.com/proemergotech/log/v3/jaegerlog"
	"github.com/uber/jaeger-client-go"
	jconfig "github.com/uber/jaeger-client-go/config"

	//%:{{ `
	"github.com/proemergotech/project-skeleton/app/config"
	"github.com/proemergotech/project-skeleton/app/rest"
	"github.com/proemergotech/project-skeleton/app/service"
	"github.com/proemergotech/project-skeleton/app/storage"
	"github.com/proemergotech/project-skeleton/app/validation"
	//%: ` | replace "project-skeleton" .ProjectName }}
)

type Container struct {
	RestServer *rest.Server
	//%: {{- if .PublicRest }}
	PublicRestServer *rest.Server
	//%: {{- end }}
	//%: {{- if .RedisCache }}
	redisCacheCloser io.Closer
	//%: {{- end }}
	//%: {{- if .RedisStore }}
	redisStoreCloser io.Closer
	//%: {{- end }}
	//%: {{- if .RedisNotice }}
	redisNoticeCloser io.Closer
	//%: {{- end }}
	traceCloser io.Closer
	//%: {{- if .Elastic }}
	elasticClient *elastic.Client
	//%: {{- end }}
}

func NewContainer(cfg *config.Config) (*Container, error) {
	c := &Container{}

	//%: {{ if .Elastic }}
	e, err := newElastic(cfg)
	if err != nil {
		return nil, err
	}
	c.elasticClient = e.ElasticClient
	//%: {{ end }}

	closer, err := newTracer(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize Jaeger Tracer")
	}
	c.traceCloser = closer

	//%: {{ if .RedisCache }}
	redisCache, err := newRedisCache(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize redis client")
	}
	c.redisCacheCloser = redisCache
	//%: {{ end }}

	//%: {{ if .RedisStore }}
	redisStore, err := newRedisStore(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize redis client")
	}
	c.redisStoreCloser = redisStore
	//%: {{ end }}

	//%: {{ if .RedisNotice }}
	redisNotice, err := newRedisNotice(cfg)
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize redis client")
	}
	c.redisNoticeCloser = redisNotice
	//%: {{ end }}

	v, err := NewValidator()
	if err != nil {
		return nil, errors.Wrap(err, "cannot initialize validator")
	}

	echoEngine := newEcho(cfg.Port, v, rest.HTTPErrorHandler)
	//%: {{ if .PublicRest }}
	publicEchoEngine := newEcho(cfg.PublicPort, v, rest.PublicHTTPErrorHandler)
	//%: {{ end }}

	svc := service.NewService()

	c.RestServer = rest.NewServer(
		echoEngine,
		rest.NewController(
			echoEngine,
			svc,
			cfg.DebugPProf,
		),
	)

	//%: {{ if .PublicRest }}
	c.PublicRestServer = rest.NewServer(
		publicEchoEngine,
		rest.NewPublicController(
			publicEchoEngine,
			svc,
		),
	)
	//%: {{ end }}

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

//%: {{ if .Elastic }}
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
} //%: {{ end }}

//%: {{ if .RedisCache }}
func newRedisCache(cfg *config.Config) (*storage.RedisCache, error) {
	redisPool, err := newRedisPool(
		cfg.RedisCachePoolIdleTimeout,
		cfg.RedisCachePoolMaxIdle,
		cfg.RedisCacheHost,
		cfg.RedisCachePort,
		cfg.RedisCacheDatabase,
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

	return storage.NewRedisCache(redisPool, redisJSON), nil
} //%: {{ end }}

//%: {{ if .RedisStore }}
func newRedisStore(cfg *config.Config) (*storage.RedisStore, error) {
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

	return storage.NewRedisStore(redisPool, redisJSON), nil
} //%: {{ end }}

//%: {{ if .RedisNotice }}
func newRedisNotice(cfg *config.Config) (*storage.RedisNotice, error) {
	redisPool, err := newRedisPool(
		cfg.RedisNoticePoolIdleTimeout,
		cfg.RedisNoticePoolMaxIdle,
		cfg.RedisNoticeHost,
		cfg.RedisNoticePort,
		cfg.RedisNoticeDatabase,
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

	return storage.NewRedisNotice(redisPool, redisJSON), nil
} //%: {{ end }}

//%: {{ if or .RedisCache .RedisStore .RedisNotice  }}
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
} //%: {{ end }}

func NewValidator() (*validation.Validator, error) {
	v := validator.New()

	v.RegisterTagNameFunc(func(field reflect.StructField) string {
		tags := []string{"param", "json", "query"}
		for _, t := range tags {
			name := strings.SplitN(field.Tag.Get(t), ",", 2)[0]
			if name != "" && name != "-" {
				return name
			}
		}
		return ""
	})

	// use it for fields with type slice and map - for these `required` isn't working as expected
	err := v.RegisterValidation("notblank", validators.NotBlank)
	if err != nil {
		return nil, err
	}

	//%: {{ if .Examples }}
	// todo: remove
	//  example validation for enums:
	// err = v.RegisterValidation("enum_status", func(fl validator.FieldLevel) bool {
	//	return schema.Statuses[fl.Field().String()]
	// })
	// if err != nil {
	//	return nil, err
	// }
	//%: {{ end }}

	return validation.NewValidator(v), nil
}

func newEcho(port int, validator *validation.Validator, httpErrorHandler echo.HTTPErrorHandler) *echo.Echo {
	e := echo.New()

	e.Use(echolog.RecoveryMiddleware(log.GlobalLogger()))
	e.HTTPErrorHandler = httpErrorHandler
	e.Validator = validator
	e.HideBanner = true
	e.HidePort = true

	e.Server = &http.Server{
		Addr:    ":" + strconv.Itoa(port),
		Handler: e,
	}

	return e
}

func (c *Container) Close() {
	//%: {{- if .Elastic }}
	c.elasticClient.Stop()
	//%: {{- end }}

	//%: {{- if .RedisCache }}
	if err := c.redisCacheCloser.Close(); err != nil {
		err = errors.Wrap(err, "redis graceful close failed")
		log.Warn(context.Background(), err.Error(), "error", err)
	}
	//%: {{- end }}

	//%: {{- if .RedisStore }}
	if err := c.redisStoreCloser.Close(); err != nil {
		err = errors.Wrap(err, "redis graceful close failed")
		log.Warn(context.Background(), err.Error(), "error", err)
	}
	//%: {{- end }}

	//%: {{- if .RedisNotice }}
	if err := c.redisNoticeCloser.Close(); err != nil {
		err = errors.Wrap(err, "redis graceful close failed")
		log.Warn(context.Background(), err.Error(), "error", err)
	}
	//%: {{- end }}

	if err := c.traceCloser.Close(); err != nil {
		err = errors.Wrap(err, "tracer graceful close failed")
		log.Warn(context.Background(), err.Error(), "error", err)
	}
}
