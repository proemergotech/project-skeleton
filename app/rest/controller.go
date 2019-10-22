package rest

import (
	"net/http"

	"github.com/labstack/echo"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/service"
	log "gitlab.com/proemergotech/log-go"
	"gitlab.com/proemergotech/log-go/echolog"
	"gitlab.com/proemergotech/trace-go/echotrace"
)

type Controller struct {
	echoEngine *echo.Echo
	svc        *service.Service
}

func NewController(
	echoEngine *echo.Echo,
	svc *service.Service,
) *Controller {
	return &Controller{
		echoEngine: echoEngine,
		svc:        svc,
	}
}

func (c *Controller) start() {

	c.echoEngine.Add(http.MethodGet, "/healthcheck", func(eCtx echo.Context) error {
		return eCtx.String(http.StatusOK, "ok")
	})

	apiRoutes := c.echoEngine.Group("/api/v1")
	apiRoutes.Use(echolog.DebugMiddleware(log.GlobalLogger(), true, true))
	apiRoutes.Use(echotrace.Middleware(opentracing.GlobalTracer(), log.GlobalLogger()))
}
