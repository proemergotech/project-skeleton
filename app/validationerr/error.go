package validationerr

import (
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/errorsf"
)

type ValidationError struct {
	Err       error
	Msg       string
	PathParam string
}

func (e ValidationError) E() error {
	msg := e.Msg
	if msg == "" {
		if e.PathParam != "" {
			msg = "invalid path parameter: " + e.PathParam
		} else {
			msg = "validation error"
		}
	}

	err := e.Err
	var details []map[string]interface{}
	if err == nil {
		err = errors.New(msg)
	} else if errs, ok := err.(validator.ValidationErrors); ok {
		err = errors.New(msg)
		details = make([]map[string]interface{}, 0, len(errs))
		for _, err := range errs {
			details = append(details, map[string]interface{}{
				"field":   err.Field(),
				"error":   service.ErrValidation,
				"message": "Field " + err.Namespace() + " failed on: " + err.Tag(),
			})
		}
	} else {
		err = errors.Wrap(err, msg)
	}

	return errorsf.WithFields(
		err,
		schema.ErrCode, service.ErrValidation,
		schema.ErrHTTPCode, 400,
		schema.ErrDetails, details,
	)
}
