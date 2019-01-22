package rest

import (
	"net/http"

	"github.com/labstack/echo"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/apierr"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
)

func DLiveRHTTPErrorHandler(err error, eCtx echo.Context) {
	if eErr, ok := err.(*echo.HTTPError); ok {
		sc := eErr.Code

		switch sc {
		case http.StatusNotFound:
			_ = eCtx.JSON(sc, &schema.HTTPError{
				Error: &schema.Error{
					Message: "route cannot be found",
					Code:    service.ErrRouteNotFound,
				}},
			)

		case http.StatusMethodNotAllowed:
			_ = eCtx.JSON(sc, &schema.HTTPError{
				Error: &schema.Error{
					Message: "method not allowed",
					Code:    service.ErrMethodNotAllowed,
				}},
			)

		default:
			_ = eCtx.JSON(http.StatusInternalServerError, &schema.HTTPError{
				Error: &schema.Error{
					Message: eErr.Error(),
					Code:    service.ErrSemanticError,
				}},
			)
		}

		return
	}

	sc := apierr.StatusCode(err)
	if sc <= 599 && sc >= 400 {
		_ = eCtx.JSON(sc, &schema.HTTPError{
			Error: &schema.Error{
				Message: err.Error(),
				Code:    apierr.Code(err),
				Details: apierr.Details(err),
			},
		})
		return
	}

	_ = eCtx.JSON(http.StatusInternalServerError, &schema.HTTPError{
		Error: &schema.Error{
			Message: "semantic",
			Code:    service.ErrSemanticError,
		}},
	)
}
