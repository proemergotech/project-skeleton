package validation

import (
	"github.com/go-playground/validator/v10"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
	"gitlab.com/proemergotech/errors"
)

type Error struct {
	Err    error
	Msg    string
	Fields []ErrorField
}

type ErrorField struct {
	Field         string
	ValidationTag string
}

func (e Error) E() error {
	msg := "validation error"
	if e.Msg != "" {
		msg += ": " + e.Msg
	}

	err := e.Err
	var details []map[string]interface{}
	var publicDetails []map[string]interface{}
	if err == nil {
		err = errors.New(msg)
	} else if errs, ok := err.(validator.ValidationErrors); ok {
		err = errors.New(msg)
		details = make([]map[string]interface{}, 0, len(errs))
		publicDetails = make([]map[string]interface{}, 0, len(errs))
		for _, err := range errs {
			details = append(details, map[string]interface{}{
				"field":   err.Field(),
				"error":   service.ErrValidation,
				"message": "Field " + err.StructNamespace() + " failed on: " + err.Tag(),
			})
			publicDetails = append(publicDetails, map[string]interface{}{
				"field":     err.Field(),
				"validator": err.Tag(),
			})
		}
	} else {
		err = errors.Wrap(err, msg)
	}

	for _, f := range e.Fields {
		details = append(details, map[string]interface{}{
			"field":   f.Field,
			"error":   service.ErrValidation,
			"message": "Field " + f.Field + " failed on: " + f.ValidationTag,
		})
		publicDetails = append(publicDetails, map[string]interface{}{
			"field":     f.Field,
			"validator": f.ValidationTag,
		})
	}

	return errors.WithFields(
		err,
		schema.ErrCode, service.ErrValidation,
		schema.ErrHTTPCode, 400,
		schema.ErrDetails, details,
		schema.ErrPublicDetails, publicDetails,
	)
}
