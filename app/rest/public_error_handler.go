package rest

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"gitlab.com/proemergotech/log-go/v3"

	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
)

func PublicDLiveRHTTPErrorHandler(err error, eCtx echo.Context) {
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
			err = service.SemanticError{Err: eErr, Fields: []interface{}{"method", eCtx.Request().Method, "url", eCtx.Request().URL.String()}}.E()
		}
	}

	httpErr, statusCode := schema.ToPublicHTTPError(err)
	_ = eCtx.JSON(statusCode, httpErr)
}
