package rest

import (
	"net/http"

	"github.com/labstack/echo"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
	"gitlab.com/proemergotech/log-go"
)

func DLiveRHTTPErrorHandler(err error, eCtx echo.Context) {
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
			err = routeNotFoundError{Err: eErr}.E()
		case http.StatusMethodNotAllowed:
			err = methodNotAllowedError{Err: eErr}.E()
		default:
			err = service.SemanticError{Err: eErr}.E()
		}
	}

	httpErr, statusCode := schema.ToHTTPError(err)
	_ = eCtx.JSON(statusCode, httpErr)
}
