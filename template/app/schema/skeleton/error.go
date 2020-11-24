//%: {{ `
package skeleton //%: ` | replace "skeleton" .SchemaPackage }}

import (
	"github.com/proemergotech/errors"

	//%:{{ `
	"github.com/proemergotech/project-skeleton/app/schema"
	//%: ` | replace "project-skeleton" .ProjectName }}
)

const (
	// 400
	ErrValidation = "ERR_VALIDATION"

	// 404
	ErrRouteNotFound = "ERR_ROUTE_NOT_FOUND"

	// 405
	ErrMethodNotAllowed = "ERR_METHOD_NOT_ALLOWED"

	// 500
	ErrClient = "ERR_CLIENT"
	//%: {{- if .Elastic }}
	ErrElastic = "ERR_ELASTIC"
	//%: {{- end }}
	//%: {{- if or .RedisCache .RedisStore .RedisNotice }}
	ErrRedis = "ERR_REDIS"
	//%: {{- end }}
	ErrSemanticError = "ERR_SEMANTIC"
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
