package rest

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"gitlab.com/proemergotech/log-go/v3"
	"gitlab.com/proemergotech/log-go/v3/echolog"
	"gitlab.com/proemergotech/trace-go/v2"
	"gitlab.com/proemergotech/trace-go/v2/echotrace"

	"gitlab.com/proemergotech/dliver-project-skeleton/app/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/validation"
)

type publicController struct {
	echoEngine *echo.Echo
	svc        *service.Service
}

func NewPublicController(echoEngine *echo.Echo, svc *service.Service) Controller {
	return &publicController{
		echoEngine: echoEngine,
		svc:        svc,
	}
}

func (pc *publicController) Start() {
	apiRoutes := pc.echoEngine.Group("/api/v1")
	startTrace, err := echotrace.Trace(trace.Start)
	if err != nil {
		log.Panic(context.Background(), err.Error(), "error", err)
	}
	apiRoutes.Use(echolog.DebugMiddleware(log.GlobalLogger(), false, false))
	apiRoutes.Use(echotrace.Middleware(opentracing.GlobalTracer(), log.GlobalLogger(), startTrace, echotrace.GenerateCorrelation(trace.NewCorrelation)))

	// todo: remove
	//  Example root
	apiRoutes.POST("/dummy", func(eCtx echo.Context) error {
		req := &struct {
			DummyData1 string `json:"dummy_data_1" validate:"required"`
			DummyData2 string `json:"dummy_data_2"`
		}{}

		if err := eCtx.Bind(req); err != nil {
			return validation.Error{Err: err, Msg: "cannot bind request"}.E()
		}

		if err := eCtx.Validate(req); err != nil {
			return err
		}

		pc.svc.SendCentrifuge(eCtx.Request().Context(), "", req)

		return eCtx.NoContent(http.StatusOK)
	})
}
