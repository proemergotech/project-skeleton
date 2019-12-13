package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
	log "gitlab.com/proemergotech/log-go"
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
			err = routeNotFoundError{Err: eErr, URL: eCtx.Request().URL.String()}.E()
		case http.StatusMethodNotAllowed:
			err = methodNotAllowedError{Err: eErr, URL: eCtx.Request().URL.String()}.E()
		default:
			err = service.SemanticError{Err: eErr, Fields: []interface{}{"url", eCtx.Request().URL.String()}}.E()
		}
	}

	httpErr, statusCode := schema.ToHTTPError(err)
	_ = eCtx.JSON(statusCode, httpErr)
}
