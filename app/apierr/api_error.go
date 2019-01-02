package apierr

import (
	"gitlab.com/proemergotech/dliver-project-skeleton/errorsf"
)

// Error is a helper type to distinguish already structured errors from non-structured ones
type Error error

func newError(cause error, statusCode int, code string, details ...ErrorDetail) Error {
	return errorsf.WithFields(
		cause,
		"status_code", statusCode,
		"code", code,
		"details", details,
	)
}

func StatusCode(e error) int {
	fI := errorsf.Field(e, "status_code")
	if fI == nil {
		return 0
	}

	return fI.(int)
}

func Code(e error) string {
	cI := errorsf.Field(e, "code")
	if cI == nil {
		return ""
	}

	return cI.(string)
}

func Details(e error) []map[string]interface{} {
	dI := errorsf.Field(e, "details")
	if dI == nil {
		return []map[string]interface{}{}
	}

	details, ok := dI.([]ErrorDetail)
	if !ok {
		return []map[string]interface{}{}
	}

	results := make([]map[string]interface{}, 0, len(details))
	for _, d := range details {
		results = append(results, d.ToMap())
	}

	return results
}
