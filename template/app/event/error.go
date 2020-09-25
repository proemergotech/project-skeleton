package event

import (
	"gitlab.com/proemergotech/errors"

	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
)

type invalidDummyEventPayloadError struct {
	Err error
}

func (e invalidDummyEventPayloadError) E() error {
	return errors.WithFields(
		errors.WrapOrNew(e.Err, "invalid dummy event body"),
		schema.ErrCode, service.ErrDummyInvalidEventPayload,
		schema.ErrHTTPCode, 400,
	)
}
