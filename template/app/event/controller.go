package event

import (
	"gitlab.com/proemergotech/geb-client-go/v2/geb"

	serviceSchema "gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/validation"
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

	// todo:
	//  add gebQueue.OnEvent handlers here

	// todo: remove
	//   Example event handler
	err := c.gebQueue.OnEvent("event/service_name/dummy/created/v1").
		Listen(func(e *geb.Event) error {
			req := &serviceSchema.AddedEvent{}

			err := e.Unmarshal(req)
			if err != nil {
				return invalidDummyEventPayloadError{Err: err}.E()
			}
			err = c.validator.Validate(req)
			if err != nil {
				return err
			}

			// handle event
			return nil
		})
	if err != nil {
		return err
	}

	return c.gebQueue.Start()
}
