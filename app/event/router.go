package event

import (
	"github.com/go-playground/validator"
	"gitlab.com/proemergotech/geb-client-go/geb"
)

type Router struct {
	gebQueue *geb.Queue
	validate *validator.Validate
}

func NewRouter(
	gebQueue *geb.Queue,
	validate *validator.Validate,
) *Router {
	return &Router{
		gebQueue: gebQueue,
		validate: validate,
	}
}

func (r *Router) route() {
}
