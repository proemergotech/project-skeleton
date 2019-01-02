package rest

import (
	"net/http"

	"gitlab.com/proemergotech/log-go"

	"gitlab.com/proemergotech/dliver-project-skeleton/app/action"

	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"

	"gitlab.com/proemergotech/trace-go/echotrace"
)

type Router struct {
	echoEngine *echo.Echo
	actions    *action.Actions
	tracer     opentracing.Tracer
}

func NewRouter(
	echoEngine *echo.Echo,
	actions *action.Actions,
	tracer opentracing.Tracer,
) *Router {
	return &Router{
		echoEngine: echoEngine,
		actions:    actions,
		tracer:     tracer,
	}
}

func (r *Router) route() {
	r.echoEngine.Add(http.MethodGet, "/healthcheck", func(eCtx echo.Context) error {
		return eCtx.String(http.StatusOK, "ok")
	})

	apiRoutes := r.echoEngine.Group("/api/v1")
	apiRoutes.Use(echotrace.Middleware(r.tracer, log.GlobalLogger()))
}
