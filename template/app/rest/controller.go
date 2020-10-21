package rest

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/pprof"
	"runtime"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/proemergotech/log-go/v3"
	"gitlab.com/proemergotech/log-go/v3/echolog"
	"gitlab.com/proemergotech/trace-go/v2/echotrace"

	//%:{{ `
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/skeleton"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/validation"
	//%: ` | replace "dliver-project-skeleton" .ProjectName | replace "skeleton" .SchemaPackage }}
)

type controller struct {
	echoEngine *echo.Echo
	svc        *service.Service
	debugPProf bool
}

func NewController(
	echoEngine *echo.Echo,
	svc *service.Service,
	debugPProf bool,
) Controller {
	return &controller{
		echoEngine: echoEngine,
		svc:        svc,
		debugPProf: debugPProf,
	}
}

func (c *controller) Start() {
	if c.debugPProf {
		runtime.SetBlockProfileRate(1)
		runtime.SetMutexProfileFraction(5)

		c.echoEngine.GET("/debug/pprof/*", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
		c.echoEngine.GET("/debug/pprof/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
		c.echoEngine.GET("/debug/pprof/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
		c.echoEngine.GET("/debug/pprof/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
		c.echoEngine.GET("/debug/pprof/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
	}

	c.echoEngine.GET("/healthcheck", func(eCtx echo.Context) error {
		return eCtx.String(http.StatusOK, "ok")
	})

	c.echoEngine.GET("metrics", echo.WrapHandler(promhttp.Handler()))

	apiRoutes := c.echoEngine.Group("/api/v1")
	apiRoutes.Use(echolog.DebugMiddleware(log.GlobalLogger(), true, true))
	apiRoutes.Use(echotrace.Middleware(opentracing.GlobalTracer(), log.GlobalLogger()))

	//%: {{ if .Examples }}
	// todo: remove
	//  Example root
	apiRoutes.POST("/dummy/:dummy_param_1", func(eCtx echo.Context) error {
		//%:{{ `
		req := &skeleton.DummyRequest{}
		//%: ` | replace "skeleton" .SchemaPackage | trim }}

		if err := eCtx.Bind(req); err != nil {
			return validation.Error{Err: err, Msg: "cannot bind request"}.E()
		}

		if err := eCtx.Validate(req); err != nil {
			return err
		}

		err := c.svc.Dummy(eCtx.Request().Context(), req)
		if err != nil {
			return err
		}

		return eCtx.NoContent(http.StatusOK)
	})
	//%: {{ end }}

	//%: {{- if and .Yafuds .Examples }}
	apiRoutes.PATCH("/dummy/:dummy_uuid", func(eCtx echo.Context) error {
		//%:{{ `
		req := &skeleton.UpdateDummyRequest{}
		//%: ` | replace "skeleton" .SchemaPackage | trim }}

		body, err := ioutil.ReadAll(eCtx.Request().Body)
		if err != nil {
			return err
		}

		reader := bytes.NewReader(body)
		eCtx.Request().Body = ioutil.NopCloser(reader)

		if err := eCtx.Bind(req); err != nil {
			return validation.Error{Err: err, Msg: "cannot bind request"}.E()
		}

		if err := eCtx.Validate(req); err != nil {
			return err
		}

		if _, err := reader.Seek(0, io.SeekStart); err != nil {
			//%:{{ `
			return skeleton.SemanticError{Err: err}.E()
			//%: ` | replace "skeleton" .SchemaPackage | trim }}
		}

		keys := make(map[string]interface{})
		if err := eCtx.Bind(&keys); err != nil {
			return validation.Error{Err: err, Msg: "cannot bind request"}.E()
		}

		resp, err := c.svc.Update(eCtx.Request().Context(), req, keys)
		if err != nil {
			return err
		}

		return eCtx.JSON(http.StatusOK, resp)
	})
	//%: {{ end }}
}
