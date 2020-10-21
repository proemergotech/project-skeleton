//%: {{ if .Geb }}
package event

import (
	"gitlab.com/proemergotech/errors"

	//%:{{ `
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/skeleton"
	//%: ` | replace "dliver-project-skeleton" .ProjectName | replace "skeleton" .SchemaPackage }}
)

type invalidDummyEventPayloadError struct {
	Err error
}

func (e invalidDummyEventPayloadError) E() error {
	return errors.WithFields(
		errors.WrapOrNew(e.Err, "invalid dummy event body"),
		//%:{{ `
		schema.ErrCode, skeleton.ErrDummyInvalidEventPayload,
		//%: ` | replace "skeleton" .SchemaPackage | trim }}
		schema.ErrHTTPCode, 400,
	)
}

//%: {{ end }}
