package service

import (
	"github.com/pkg/errors"
	"gitlab.com/proemergotech/centrifuge-client-go/api"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/errorsf"
)

type centrifugeError struct {
	Err  error
	CErr *api.Error
}

func (e centrifugeError) E() error {
	err := e.Err
	var details []map[string]interface{}
	if e.CErr != nil {
		err = errors.New("centrifuge replied with error: " + e.CErr.Message)
		details = []map[string]interface{}{
			{
				"code":    e.CErr.Code,
				"message": e.CErr.Message,
			},
		}
	} else if err == nil {
		err = errors.New("centrifuge error")
	} else {
		err = errors.Wrap(err, "centrifuge error")
	}

	return errorsf.WithFields(
		err,
		schema.ErrCode, service.ErrCentrifuge,
		schema.ErrHTTPCode, 500,
		schema.ErrDetails, details,
	)
}
