package apierr

import "gitlab.com/proemergotech/dliver-project-skeleton/app/client/centrifugo/proto/apiproto"

type ErrorDetail interface {
	ToMap() map[string]interface{}
}

type ValidationErrorDetail struct {
	Field   string
	Error   string
	Message string
}

func (ved *ValidationErrorDetail) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"field":   ved.Field,
		"error":   ved.Error,
		"message": ved.Message,
	}
}

type CentrifugeErrorDetail struct {
	Error *apiproto.Error
}

func (ced *CentrifugeErrorDetail) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"code":    ced.Error.Code,
		"message": ced.Error.Message,
	}
}
