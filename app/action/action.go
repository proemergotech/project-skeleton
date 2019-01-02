package action

import (
	"github.com/go-playground/validator"
	"gitlab.com/proemergotech/dliver-project-skeleton/app"
)

type Actions struct {
	core     *app.Core
	validate *validator.Validate
}

func NewActions(
	core *app.Core,
	validate *validator.Validate,
) *Actions {
	return &Actions{
		core:     core,
		validate: validate,
	}
}
