package event

import (
	"github.com/go-playground/validator"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/service"
	"gitlab.com/proemergotech/geb-client-go/geb"
)

type Controller struct {
	gebQueue *geb.Queue
	validate *validator.Validate
	svc      *service.Service
}

func NewController(
	gebQueue *geb.Queue,
	validate *validator.Validate,
	svc *service.Service,
) *Controller {
	return &Controller{
		gebQueue: gebQueue,
		validate: validate,
		svc:      svc,
	}
}

func (c *Controller) start() error {

	// TODO:
	//  add gebQueue.OnEvent handlers here

	c.gebQueue.Start()

	return nil
}
