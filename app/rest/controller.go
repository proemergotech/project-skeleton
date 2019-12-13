package rest

import (
	"net/http"
	"net/http/pprof"
	"runtime"

	"github.com/labstack/echo/v4"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/validation"
	log "gitlab.com/proemergotech/log-go"
	"gitlab.com/proemergotech/log-go/echolog"
	"gitlab.com/proemergotech/trace-go/echotrace"
)

type Controller struct {
	echoEngine *echo.Echo
	svc        *service.Service
	debugPProf bool
}

func NewController(
	echoEngine *echo.Echo,
	svc *service.Service,
	debugPProf bool,
) *Controller {
	return &Controller{
		echoEngine: echoEngine,
		svc:        svc,
		debugPProf: debugPProf,
	}
}

func (c *Controller) start() {
	if c.debugPProf {
		runtime.SetBlockProfileRate(1)
		runtime.SetMutexProfileFraction(5)

		c.echoEngine.GET("/debug/pprof/*", echo.WrapHandler(http.HandlerFunc(pprof.Index)))
		c.echoEngine.GET("/debug/pprof/cmdline", echo.WrapHandler(http.HandlerFunc(pprof.Cmdline)))
		c.echoEngine.GET("/debug/pprof/profile", echo.WrapHandler(http.HandlerFunc(pprof.Profile)))
		c.echoEngine.GET("/debug/pprof/symbol", echo.WrapHandler(http.HandlerFunc(pprof.Symbol)))
		c.echoEngine.GET("/debug/pprof/trace", echo.WrapHandler(http.HandlerFunc(pprof.Trace)))
	}

	c.echoEngine.Add(http.MethodGet, "/healthcheck", func(eCtx echo.Context) error {
		return eCtx.String(http.StatusOK, "ok")
	})

	c.echoEngine.GET("metrics", echo.WrapHandler(promhttp.Handler()))

	apiRoutes := c.echoEngine.Group("/api/v1")
	apiRoutes.Use(echolog.DebugMiddleware(log.GlobalLogger(), true, true))
	apiRoutes.Use(echotrace.Middleware(opentracing.GlobalTracer(), log.GlobalLogger()))

	// todo: remove
	//  Example root
	apiRoutes.Add(http.MethodPost, "/dummy", func(eCtx echo.Context) error {
		req := &struct {
			DummyData1 string `json:"dummy_data_1"`
			DummyData2 string `json:"dummy_data_2"`
		}{}

		if err := eCtx.Bind(req); err != nil {
			return validation.Error{Err: err, Msg: "cannot bind request"}.E()
		}

		if err := eCtx.Validate(req); err != nil {
			return err
		}

		c.svc.SendCentrifuge(eCtx.Request().Context(), "", "", req)

		return eCtx.NoContent(http.StatusOK)
	})
}
