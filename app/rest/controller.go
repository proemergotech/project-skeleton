package rest

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/service"
	"gitlab.com/proemergotech/log-go"
	"gitlab.com/proemergotech/trace-go/echotrace"
)

type Controller struct {
	echoEngine *echo.Echo
	service    *service.Service
	tracer     opentracing.Tracer
}

func NewController(
	echoEngine *echo.Echo,
	service *service.Service,
	tracer opentracing.Tracer,
) *Controller {
	return &Controller{
		echoEngine: echoEngine,
		service:    service,
		tracer:     tracer,
	}
}

func (c *Controller) start() {
	c.echoEngine.Add(http.MethodGet, "/healthcheck", func(eCtx echo.Context) error {
		return eCtx.String(http.StatusOK, "ok")
	})

	apiRoutes := c.echoEngine.Group("/api/v1")
	apiRoutes.Use(echotrace.Middleware(c.tracer, log.GlobalLogger()))
}
