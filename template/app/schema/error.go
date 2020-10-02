package schema

import (
	"gitlab.com/proemergotech/errors"
)

const (
	ErrCode    = "code"
	ErrDetails = "details"
	//%: {{ if .PublicRest }}
	ErrPublicDetails = "public_details"
	//%: {{ end }}
	ErrHTTPCode = "http_code"
)

type HTTPError struct {
	Error Error `json:"error"`
}

type Error struct {
	Message string                   `json:"message,omitempty" validate:"omitempty"`
	Code    string                   `json:"code" validate:"required"`
	Details []map[string]interface{} `json:"details,omitempty" validate:"omitempty"`
}

func ToHTTPError(err error) (*HTTPError, int) {
	httpCode := ErrorHTTPCode(err)
	if httpCode == 0 {
		httpCode = 500
	}

	return &HTTPError{
		Error: Error{
			Message: err.Error(),
			Code:    ErrorCode(err),
			Details: ErrorDetails(err),
		},
	}, httpCode
}

//%: {{ if .PublicRest }}
func ToPublicHTTPError(err error) (*HTTPError, int) {
	httpCode := ErrorHTTPCode(err)
	if httpCode == 0 {
		httpCode = 500
	}

	code := ErrorCode(err)
	if httpCode >= 500 {
		code = "ERR_INTERNAL"
	}

	return &HTTPError{
		Error: Error{
			Code:    code,
			Details: ErrorPublicDetails(err),
		},
	}, httpCode
} //%: {{ end }}

func ErrorCode(err error) string {
	field := errors.Field(err, ErrCode)
	if field == nil {
		return ""
	}

	return field.(string)
}

func ErrorHTTPCode(err error) int {
	field := errors.Field(err, ErrHTTPCode)
	if field == nil {
		return 0
	}

	return field.(int)
}

func ErrorDetails(err error) []map[string]interface{} {
	field := errors.Field(err, ErrDetails)
	if field == nil {
		return []map[string]interface{}{}
	}

	return field.([]map[string]interface{})
}

//%: {{ if .PublicRest }}
func ErrorPublicDetails(err error) []map[string]interface{} {
	field := errors.Field(err, ErrPublicDetails)
	if field == nil {
		return []map[string]interface{}{}
	}

	return field.([]map[string]interface{})
} //%: {{ end }}
