package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/proemergotech/errors"

	//%:{{ `
	"github.com/proemergotech/project-skeleton/app/schema"
	"github.com/proemergotech/project-skeleton/app/schema/skeleton"
	//%: ` | replace "project-skeleton" .ProjectName | replace "skeleton" .SchemaPackage }}
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
	//%: {{ if .PublicRest }}
	var publicDetails []map[string]interface{}
	//%: {{ end }}
	if err == nil {
		err = errors.New(msg)
	} else if errs, ok := err.(validator.ValidationErrors); ok {
		err = errors.New(msg)
		details = make([]map[string]interface{}, 0, len(errs))
		//%: {{ if .PublicRest }}
		publicDetails = make([]map[string]interface{}, 0, len(errs))
		//%: {{ end }}
		for _, err := range errs {
			details = append(details, map[string]interface{}{
				"field": err.Field(),
				//%:{{ `
				"error": skeleton.ErrValidation,
				//%: ` | replace "skeleton" .SchemaPackage }}
				"message": "Field " + err.StructNamespace() + " failed on: " + err.Tag(),
			})
			//%: {{ if .PublicRest }}
			publicDetails = append(publicDetails, map[string]interface{}{
				"field":     err.Field(),
				"validator": err.Tag(),
			})
			//%: {{ end }}
		}
	} else {
		err = errors.Wrap(err, msg)
	}

	for _, f := range e.Fields {
		details = append(details, map[string]interface{}{
			"field": f.Field,
			//%:{{ `
			"error": skeleton.ErrValidation,
			//%: ` | replace "skeleton" .SchemaPackage }}
			"message": "Field " + f.Field + " failed on: " + f.ValidationTag,
		})
		//%: {{ if .PublicRest }}
		publicDetails = append(publicDetails, map[string]interface{}{
			"field":     f.Field,
			"validator": f.ValidationTag,
		})
		//%: {{ end }}
	}

	return errors.WithFields(
		err,
		//%:{{ `
		schema.ErrCode, skeleton.ErrValidation,
		//%: ` | replace "skeleton" .SchemaPackage }}
		schema.ErrHTTPCode, 400,
		schema.ErrDetails, details,
		//%: {{ if .PublicRest }}
		schema.ErrPublicDetails, publicDetails,
		//%: {{ end }}
	)
}
