//%: {{- if or .Elastic .RedisCache .RedisStore .RedisNotice .Yafuds }}
package storage

import (
	"gitlab.com/proemergotech/errors"

	//%:{{ `
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/skeleton"
	//%: ` | replace "dliver-project-skeleton" .ProjectName | replace "skeleton" .SchemaPackage }}
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
		//%: ` | replace "skeleton" .SchemaPackage }}
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
		//%: ` | replace "skeleton" .SchemaPackage }}
		schema.ErrHTTPCode, 500,
	)
} //%: {{- end }}

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

	return errors.WithFields(
		errors.WrapOrNew(e.Err, msg),
		//%:{{ `
		schema.ErrCode, skeleton.ErrYafuds,
		//%: ` | replace "skeleton" .SchemaPackage }}
		schema.ErrHTTPCode, 500,
	)
} //%: {{- end }}

//%: {{- end }}
