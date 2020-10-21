//%: {{ if or .Yafuds .Geb .Elastic  }}
package service

import (
	"gitlab.com/proemergotech/errors"

	//%:{{ `
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/skeleton"
	//%: ` | replace "dliver-project-skeleton" .ProjectName | replace "skeleton" .SchemaPackage }}
)

//%: {{- if .Yafuds }}
type yafudsError struct {
	Err error
	Msg string
}

func (e yafudsError) E() error {
	msg := "yafuds error"
	if e.Msg != "" {
		msg += ": " + e.Msg
	}

	err := e.Err
	if err == nil {
		err = errors.New(msg)
	} else {
		err = errors.Wrap(err, msg)
	}

	//%:{{ `
	return errors.WithFields(err, schema.ErrCode, skeleton.ErrYafuds, schema.ErrHTTPCode, 500)
	//%: ` | replace "skeleton" .SchemaPackage | trim }}
} //%: {{- end }}

//%: {{- if .Geb }}
type gebError struct {
	Err error
	Msg string
}

func (e gebError) E() error {
	msg := "geb error"
	if e.Msg != "" {
		msg += ": " + e.Msg
	}

	err := e.Err
	if err == nil {
		err = errors.New(msg)
	} else {
		err = errors.Wrap(err, msg)
	}

	//%:{{ `
	return errors.WithFields(err, schema.ErrCode, skeleton.ErrGeb, schema.ErrHTTPCode, 500)
	//%: ` | replace "skeleton" .SchemaPackage | trim }}
} //%: {{- end }}

//%: {{- end }}
