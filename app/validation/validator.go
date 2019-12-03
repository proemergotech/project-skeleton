package validation

import "github.com/go-playground/validator/v10"

type Validator struct {
	validator *validator.Validate
}

func NewValidator(
	validator *validator.Validate,
) *Validator {
	return &Validator{
		validator: validator,
	}
}

func (v *Validator) Validate(i interface{}) error {
	if err := v.validator.Struct(i); err != nil {
		return Error{Err: err}.E()
	}

	return nil
}
