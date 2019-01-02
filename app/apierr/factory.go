package apierr

import (
	"net/http"

	"gitlab.com/proemergotech/dliver-project-skeleton/app/client/centrifugo/proto/apiproto"

	"github.com/go-playground/validator"
	"github.com/pkg/errors"

	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
)

func Semantic(cause error) Error {
	return newError(
		errors.Wrap(cause, "semantic"),
		http.StatusInternalServerError,
		service.ErrSemanticError,
	)
}

func Validation(err error) Error {
	errs, ok := err.(validator.ValidationErrors)
	if !ok {
		return Semantic(errors.Errorf("error passed to apierr.Validation must be of type validator.ValidationErrors, %T found", err))
	}

	details := make([]ErrorDetail, 0, len(errs))
	for _, err := range errs {
		detail := &ValidationErrorDetail{
			Field:   err.Field(),
			Error:   service.ErrValidation,
			Message: "Field " + err.Namespace() + " failed on: " + err.Tag(),
		}

		details = append(details, ErrorDetail(detail))
	}

	return newError(
		errors.New("validation error"),
		http.StatusBadRequest,
		service.ErrValidation,
		details...,
	)
}

func Centrifuge(err error) Error {
	return newError(
		errors.Wrap(err, "error connecting to centrifuge"),
		http.StatusInternalServerError,
		service.ErrCentrifuge,
	)
}

func CentrifugeResponse(err *apiproto.Error) Error {
	return newError(
		errors.New("centrifuge replied with error: "+err.Message),
		http.StatusInternalServerError,
		service.ErrCentrifuge,
		&CentrifugeErrorDetail{Error: err},
	)
}

func RedisUnavailable(cause error) Error {
	return newError(
		errors.Wrap(cause, "redis error"),
		http.StatusInternalServerError,
		service.ErrRedisUnavailable,
	)
}
