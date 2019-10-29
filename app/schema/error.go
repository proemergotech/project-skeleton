package schema

import (
	"gitlab.com/proemergotech/dliver-project-skeleton/errorsf"
)

const (
	ErrCode     = "code"
	ErrDetails  = "details"
	ErrHTTPCode = "http_code"
	ErrGRPCCode = "grcp_code"
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
			Code: code,
		},
	}, httpCode
}

func ErrorCode(err error) string {
	field := errorsf.Field(err, ErrCode)
	if field == nil {
		return ""
	}

	return field.(string)
}

func ErrorHTTPCode(err error) int {
	field := errorsf.Field(err, ErrHTTPCode)
	if field == nil {
		return 0
	}

	return field.(int)
}

func ErrorGRPCCode(err error) int {
	field := errorsf.Field(err, ErrGRPCCode)
	if field == nil {
		return 0
	}

	return field.(int)
}

func ErrorDetails(err error) []map[string]interface{} {
	field := errorsf.Field(err, ErrDetails)
	if field == nil {
		return []map[string]interface{}{}
	}

	return field.([]map[string]interface{})
}
