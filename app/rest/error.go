package rest

import (
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
	"gitlab.com/proemergotech/errors"
)

type routeNotFoundError struct {
	Err error
	URL string
}

func (e routeNotFoundError) E() error {
	msg := "route cannot be found"

	if e.URL != "" {
		msg += ": '" + e.URL + "'"
	}

	return errors.WithFields(
		errors.WrapOrNew(e.Err, msg),
		schema.ErrHTTPCode, 404,
		schema.ErrCode, service.ErrRouteNotFound,
	)
}

type methodNotAllowedError struct {
	Err error
	URL string
}

func (e methodNotAllowedError) E() error {
	msg := "method not allowed"

	if e.URL != "" {
		msg += ": '" + e.URL + "'"
	}

	return errors.WithFields(
		errors.WrapOrNew(e.Err, msg),
		schema.ErrHTTPCode, 405,
		schema.ErrCode, service.ErrMethodNotAllowed,
	)
}
