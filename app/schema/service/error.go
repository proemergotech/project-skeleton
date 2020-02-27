package service

import (
	"github.com/pkg/errors"

	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/errorsf"
)

const (
	// 400
	ErrValidation = "ERR_VALIDATION"

	// 404
	ErrRouteNotFound = "ERR_ROUTE_NOT_FOUND"

	// 405
	ErrMethodNotAllowed = "ERR_METHOD_NOT_ALLOWED"

	// 500
	ErrCentrifuge    = "ERR_CENTRIFUGE"
	ErrElastic       = "ERR_ELASTIC"
	ErrRedis         = "ERR_REDIS"
	ErrSemanticError = "ERR_SEMANTIC"
	ErrYafuds        = "ERR_YAFUDS"
	ErrClient        = "ERR_CLIENT"
)

type SemanticError struct {
	Err    error
	Msg    string
	Fields []interface{}
}

func (e SemanticError) E() error {
	msg := "semantic error"
	if e.Msg != "" {
		msg += ": " + e.Msg
	}

	err := e.Err
	if err == nil {
		err = errors.New(msg)
	} else {
		err = errors.Wrap(err, msg)
	}

	return errorsf.WithFields(
		err,
		append(e.Fields,
			schema.ErrCode, ErrSemanticError,
			schema.ErrHTTPCode, 500,
		)...,
	)
}
