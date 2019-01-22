package service

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
