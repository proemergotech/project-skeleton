package storage

import (
	"github.com/pkg/errors"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/errorsf"
)

type redisUnavailableError struct {
	Err error
	Msg string
}

func (e redisUnavailableError) E() error {
	msg := e.Msg
	if msg == "" {
		msg = "redis error"
	}

	err := e.Err
	if err == nil {
		err = errors.New(msg)
	} else {
		err = errors.Wrap(err, msg)
	}

	return errorsf.WithFields(
		err,
		schema.ErrCode, service.ErrRedisUnavailable,
		schema.ErrHTTPCode, 500,
	)
}
