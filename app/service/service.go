package service

import (
	"context"
	"encoding/json"
	"fmt"

	"gitlab.com/proemergotech/yafuds-client-go/client"

	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"

	"gitlab.com/proemergotech/centrifuge-client-go/api"
	"gitlab.com/proemergotech/log-go"
)

type Service struct {
	centrifugeClient api.CentrifugeClient
	yafudsClient     *client.Client
}

func NewService(
	centrifugeClient api.CentrifugeClient,
	yafudsClient *client.Client,
) *Service {
	return &Service{
		centrifugeClient: centrifugeClient,
		yafudsClient:     yafudsClient,
	}
}

// Centrifuge example
func (s *Service) SendCentrifuge(ctx context.Context, namespace string, identifier string, eventData interface{}) {
	centrifugeChannel := namespace + ":" + identifier

	data, err := json.Marshal(eventData)

	if err != nil {
		err = service.SemanticError{Err: err, Msg: fmt.Sprintf("unable to marshal eventData of type: %T", eventData)}.E()
		log.Error(ctx, err.Error(), "error", err)
		return
	}

	resp, err := s.centrifugeClient.Publish(ctx, &api.PublishRequest{
		Channel: centrifugeChannel,
		Data:    data,
	})
	if err != nil {
		err = centrifugeError{Err: err}.E()
		log.Error(ctx, err.Error(), "error", err)
		return
	}
	if resp.Error != nil {
		err = centrifugeError{CErr: resp.Error}.E()
		log.Error(ctx, err.Error(), "error", err)
		return
	}
}
