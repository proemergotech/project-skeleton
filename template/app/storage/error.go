//%: {{- if or .Elastic .RedisCache .RedisStore .RedisNotice }}
package storage

import (
	"github.com/proemergotech/errors"

	//%:{{ `
	"github.com/proemergotech/project-skeleton/app/schema"
	"github.com/proemergotech/project-skeleton/app/schema/skeleton"
	//%: ` | replace "project-skeleton" .ProjectName | replace "skeleton" .SchemaPackage }}
)

//%: {{- if .Elastic }}
type elasticError struct {
	Err error
	Msg string
}

func (e elasticError) E() error {
	msg := "elastic error"
	if e.Msg != "" {
		msg += ": " + e.Msg
	}

	return errors.WithFields(
		errors.WrapOrNew(e.Err, msg),
		//%:{{ `
		schema.ErrCode, skeleton.ErrElastic,
		//%: ` | replace "skeleton" .SchemaPackage | trim }}
		schema.ErrHTTPCode, 500,
	)
} //%: {{- end }}

//%: {{- if or .RedisCache .RedisStore .RedisNotice }}
type redisError struct {
	Err error
	Msg string
}

func (e redisError) E() error {
	msg := "redis error"
	if e.Msg != "" {
		msg += ": " + e.Msg
	}

	return errors.WithFields(
		errors.WrapOrNew(e.Err, msg),
		//%:{{ `
		schema.ErrCode, skeleton.ErrRedis,
		//%: ` | replace "skeleton" .SchemaPackage | trim }}
		schema.ErrHTTPCode, 500,
	)
} //%: {{- end }}

//%: {{- end }}
