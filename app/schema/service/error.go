package service

const ErrFieldPrefix = "service_"

const (
	// 400
	ErrValidation = "ERR_VALIDATION"

	// 404
	ErrRouteNotFound = "ERR_ROUTE_NOT_FOUND"

	// 405
	ErrMethodNotAllowed = "ERR_METHOD_NOT_ALLOWED"

	// 500
	ErrCentrifuge       = "ERR_CENTRIFUGE"
	ErrRedisUnavailable = "ERR_REDIS_UNAVAILABLE"
	ErrSemanticError    = "ERR_SEMANTIC_ERROR"
)

type HTTPError struct {
	Error Error `json:"error"`
}

type Error struct {
	Message string                   `json:"message" validate:"required"`
	Code    string                   `json:"code" validate:"required"`
	Details []map[string]interface{} `json:"details,omitempty" validate:"omitempty"`
}

func (e *Error) Error() string {
	msg := e.Message

	return msg
}

func (e *Error) Fields() []interface{} {
	return []interface{}{
		ErrFieldPrefix + "code", e.Code,
		ErrFieldPrefix + "message", e.Message,
		ErrFieldPrefix + "details", e.Details,
	}
}
