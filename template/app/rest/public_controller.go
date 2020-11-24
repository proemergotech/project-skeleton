//%: {{ if .PublicRest }}
package rest

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/opentracing/opentracing-go"
	"github.com/proemergotech/log/v3"
	"github.com/proemergotech/log/v3/echolog"
	"github.com/proemergotech/trace/v2"
	"github.com/proemergotech/trace/v2/echotrace"

	//%:{{ `
	"github.com/proemergotech/project-skeleton/app/schema/skeleton"
	"github.com/proemergotech/project-skeleton/app/service"
	"github.com/proemergotech/project-skeleton/app/validation"
	//%: ` | replace "project-skeleton" .ProjectName | replace "skeleton" .SchemaPackage }}
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

		err := pc.svc.Dummy(eCtx.Request().Context(), req)
		if err != nil {
			return err
		}

		return eCtx.NoContent(http.StatusOK)
	})
	//%: {{ end }}
}

//%: {{ end }}
