package storage

import (
	"gitlab.com/proemergotech/errors"

	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
)

type elasticError struct {
	Err error
	Msg string
}

func (e elasticError) E() error {
	msg := "elastic error"
	if e.Msg != "" {
		msg += ": " + e.Msg
	}

	return errors.WithFields(
		errors.WrapOrNew(e.Err, msg),
		schema.ErrCode, service.ErrElastic,
		schema.ErrHTTPCode, 500,
	)
}

type redisError struct {
	Err error
	Msg string
}

func (e redisError) E() error {
	msg := "redis error"
	if e.Msg != "" {
		msg += ": " + e.Msg
	}

	return errors.WithFields(
		errors.WrapOrNew(e.Err, msg),
		schema.ErrCode, service.ErrRedis,
		schema.ErrHTTPCode, 500,
	)
}

type yafudsError struct {
	Err error
	Msg string
}

func (e yafudsError) E() error {
	msg := "yafuds error"
	if e.Msg != "" {
		msg += ": " + e.Msg
	}

	return errors.WithFields(
		errors.WrapOrNew(e.Err, msg),
		schema.ErrCode, service.ErrYafuds,
		schema.ErrHTTPCode, 500,
	)
}
