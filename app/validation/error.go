package validation

import (
	"github.com/go-playground/validator"
	"github.com/pkg/errors"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/errorsf"
)

type Error struct {
	Err       error
	Msg       string
	PathParam string
}

func (e Error) E() error {
	msg := "validation error"
	if e.PathParam != "" {
		msg += ": invalid path parameter: " + e.PathParam
	}
	if e.Msg != "" {
		msg += ": " + e.Msg
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
				"message": "Field " + err.StructNamespace() + " failed on: " + err.Tag(),
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
