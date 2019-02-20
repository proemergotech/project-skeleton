package service

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"gitlab.com/proemergotech/centrifuge-client-go/api"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/apierr"
	log "gitlab.com/proemergotech/log-go"
)

type Service struct {
	centrifugeClient api.CentrifugeClient
}

func NewService(
	centrifugeClient api.CentrifugeClient,
) *Service {
	return &Service{
		centrifugeClient: centrifugeClient,
	}
}

// Centrifuge example
func (s *Service) SendCentrifuge(ctx context.Context, namespace string, identifier string, eventData interface{}) {
	centrifugeChannel := namespace + ":" + identifier

	data, err := json.Marshal(eventData)

	if err != nil {
		err = apierr.Semantic(errors.Wrapf(err, "unable to marshal eventData of type: %T", eventData))
		log.Error(ctx, err.Error(), "error", err)
		return
	}

	resp, err := s.centrifugeClient.Publish(ctx, &api.PublishRequest{
		Channel: centrifugeChannel,
		Data:    data,
	})
	if err != nil {
		err = apierr.Centrifuge(err)
		log.Error(ctx, err.Error(), "error", err)
		return
	}
	if resp.Error != nil {
		err = apierr.CentrifugeResponse(resp.Error)
		log.Error(ctx, err.Error(), "error", err)
		return
	}
}
