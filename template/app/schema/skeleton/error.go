//%: {{ `
package skeleton //%: ` | replace "skeleton" .SchemaPackage }}

import (
	"gitlab.com/proemergotech/errors"

	//%:{{ `
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	//%: ` | replace "dliver-project-skeleton" .ProjectName }}
)

const (
	// 400
	ErrValidation = "ERR_VALIDATION"
	//%: {{- if .Geb }}
	ErrDummyInvalidEventPayload = "ERR_DUMMY_INVALID_EVENT_PAYLOAD"
	//%: {{- end }}

	// 404
	ErrRouteNotFound = "ERR_ROUTE_NOT_FOUND"

	// 405
	ErrMethodNotAllowed = "ERR_METHOD_NOT_ALLOWED"

	// 500
	//%: {{- if .Elastic }}
	ErrElastic = "ERR_ELASTIC"
	//%: {{- end }}
	//%: {{- if or .RedisCache .RedisStore .RedisNotice }}
	ErrRedis = "ERR_REDIS"
	//%: {{- end }}
	ErrSemanticError = "ERR_SEMANTIC"
	//%: {{- if .Yafuds }}
	ErrYafuds = "ERR_YAFUDS"
	//%: {{- end }}
	//%: {{- if or .Centrifuge .SiteConfig }}
	ErrClient = "ERR_CLIENT"
	//%: {{- end }}
)

type SemanticError struct {
	Err    error
	Msg    string
	Fields []interface{}
}

func (e SemanticError) E() error {
	msg := "semantic error"
	if e.Msg != "" {
		msg += ": " + e.Msg
	}

	return errors.WithFields(
		errors.WrapOrNew(e.Err, msg),
		append(e.Fields,
			schema.ErrCode, ErrSemanticError,
			schema.ErrHTTPCode, 500,
		)...,
	)
}
