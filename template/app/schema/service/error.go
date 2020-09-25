package service

import (
	"gitlab.com/proemergotech/errors"

	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
)

const (
	// 400
	ErrValidation               = "ERR_VALIDATION"
	ErrDummyInvalidEventPayload = "ERR_DUMMY_INVALID_EVENT_PAYLOAD"

	// 404
	ErrRouteNotFound = "ERR_ROUTE_NOT_FOUND"

	// 405
	ErrMethodNotAllowed = "ERR_METHOD_NOT_ALLOWED"

	// 500
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

	return errors.WithFields(
		errors.WrapOrNew(e.Err, msg),
		append(e.Fields,
			schema.ErrCode, ErrSemanticError,
			schema.ErrHTTPCode, 500,
		)...,
	)
}
