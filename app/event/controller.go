package event

import (
	"gitlab.com/proemergotech/dliver-project-skeleton/app/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/validation"
	"gitlab.com/proemergotech/geb-client-go/v2/geb"
)

type Controller struct {
	gebQueue  *geb.Queue
	validator *validation.Validator
	svc       *service.Service
}

func NewController(
	gebQueue *geb.Queue,
	validator *validation.Validator,
	svc *service.Service,
) *Controller {
	return &Controller{
		gebQueue:  gebQueue,
		validator: validator,
		svc:       svc,
	}
}

func (c *Controller) start() error {

	// TODO:
	//  add gebQueue.OnEvent handlers here

	return c.gebQueue.Start()
}
