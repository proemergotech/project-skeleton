package rest

import (
	"net/http"

	"github.com/labstack/echo"
	opentracing "github.com/opentracing/opentracing-go"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/validation"
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
