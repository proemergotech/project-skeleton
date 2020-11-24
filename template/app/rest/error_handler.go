package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/proemergotech/log/v3"

	//%:{{ `
	"github.com/proemergotech/project-skeleton/app/schema"
	"github.com/proemergotech/project-skeleton/app/schema/skeleton"
	//%: ` | replace "project-skeleton" .ProjectName | replace "skeleton" .SchemaPackage }}
)

func HTTPErrorHandler(err error, eCtx echo.Context) {
	defer func() {
		sc := eCtx.Response().Status
		if sc >= 400 && sc < 500 {
			log.Warn(eCtx.Request().Context(), err.Error(), "error", err)
		} else {
			log.Error(eCtx.Request().Context(), err.Error(), "error", err)
		}
	}()

	if eErr, ok := err.(*echo.HTTPError); ok {
		sc := eErr.Code

		switch sc {
		case http.StatusNotFound:
			err = routeNotFoundError{Err: eErr, Method: eCtx.Request().Method, URL: eCtx.Request().URL.String()}.E()
		case http.StatusMethodNotAllowed:
			err = methodNotAllowedError{Err: eErr, Method: eCtx.Request().Method, URL: eCtx.Request().URL.String()}.E()
		default:
			//%:{{ `
			err = skeleton.SemanticError{Err: eErr, Fields: []interface{}{"method", eCtx.Request().Method, "url", eCtx.Request().URL.String()}}.E()
			//%: ` | replace "skeleton" .SchemaPackage | trim }}
		}
	}

	httpErr, statusCode := schema.ToHTTPError(err)
	_ = eCtx.JSON(statusCode, httpErr)
}
