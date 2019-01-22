package schema

type HTTPError struct {
	Error *Error `json:"error"`
}

type Error struct {
	Message string                   `json:"message" validate:"required"`
	Code    string                   `json:"code" validate:"required"`
	Details []map[string]interface{} `json:"details,omitempty" validate:"omitempty"`
	prefix  string
}

func NewError(prefix string) *Error {
	return &Error{
		prefix: prefix + "_",
	}
}

func (e *Error) Error() string {
	msg := e.Message

	return msg
}

func (e *Error) Fields() []interface{} {
	return []interface{}{
		e.prefix + "code", e.Code,
		e.prefix + "message", e.Message,
		e.prefix + "details", e.Details,
	}
}
