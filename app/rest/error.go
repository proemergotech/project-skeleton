package rest

import (
	"github.com/pkg/errors"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/errorsf"
)

type routeNotFoundError struct {
	Err error
}

func (e routeNotFoundError) E() error {
	return errorsf.WithFields(
		errors.Wrap(e.Err, "route cannot be found"),
		schema.ErrHTTPCode, 404,
		schema.ErrCode, service.ErrRouteNotFound,
	)
}

type methodNotAllowedError struct {
	Err error
}

func (e methodNotAllowedError) E() error {
	return errorsf.WithFields(
		errors.Wrap(e.Err, "method not allowed"),
		schema.ErrHTTPCode, 405,
		schema.ErrCode, service.ErrMethodNotAllowed,
	)
}
