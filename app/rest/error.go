package rest

import (
	"github.com/pkg/errors"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/errorsf"
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

	err := e.Err
	if err == nil {
		err = errors.New(msg)
	} else {
		err = errors.Wrap(err, msg)
	}

	return errorsf.WithFields(
		err,
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

	err := e.Err
	if err == nil {
		err = errors.New(msg)
	} else {
		err = errors.Wrap(err, msg)
	}

	return errorsf.WithFields(
		err,
		schema.ErrHTTPCode, 405,
		schema.ErrCode, service.ErrMethodNotAllowed,
	)
}
