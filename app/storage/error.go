package storage

import (
	"github.com/pkg/errors"

	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/errorsf"
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

	err := e.Err
	if err == nil {
		err = errors.New(msg)
	} else {
		err = errors.Wrap(err, msg)
	}

	return errorsf.WithFields(err, schema.ErrCode, service.ErrElastic, schema.ErrHTTPCode, 500)
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

	err := e.Err
	if err == nil {
		err = errors.New(msg)
	} else {
		err = errors.Wrap(err, msg)
	}

	return errorsf.WithFields(
		err,
		schema.ErrCode, service.ErrRedis,
		schema.ErrHTTPCode, 500,
	)
}
