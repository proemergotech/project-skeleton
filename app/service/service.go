package service

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"gitlab.com/proemergotech/centrifuge-client-go/v2/api"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/client"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
	log "gitlab.com/proemergotech/log-go/v3"
	yafuds "gitlab.com/proemergotech/yafuds-client-go/client"
)

type Service struct {
	centrifugeClient api.CentrifugeClient
	centrifugeJSON   jsoniter.API
	yafudsClient     yafuds.Client
	dummyClient      *client.Dummy
}

func NewService(
	centrifugeClient api.CentrifugeClient,
	centrifugeJSON jsoniter.API,
	yafudsClient yafuds.Client,
	dummyClient *client.Dummy,
) *Service {
	return &Service{
		centrifugeClient: centrifugeClient,
		centrifugeJSON:   centrifugeJSON,
		yafudsClient:     yafudsClient,
		dummyClient:      dummyClient,
	}
}

// Centrifuge example
func (s *Service) SendCentrifuge(ctx context.Context, centrifugeChannel string, eventData interface{}) {
	data, err := s.centrifugeJSON.Marshal(eventData)
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
