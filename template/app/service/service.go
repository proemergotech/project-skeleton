package service

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"gitlab.com/proemergotech/log-go/v3"

	"gitlab.com/proemergotech/dliver-project-skeleton/app/client"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/centrifuge"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/storage"
)

type Service struct {
	centrifugeClient *client.Centrifuge
	centrifugeJSON   jsoniter.API
	yafudsClient     *storage.Yafuds
	siteConfigClient *client.SiteConfig
}

func NewService(
	centrifugeClient *client.Centrifuge,
	centrifugeJSON jsoniter.API,
	yafudsStorage *storage.Yafuds,
	siteConfigClient *client.SiteConfig,
) *Service {
	return &Service{
		centrifugeClient: centrifugeClient,
		centrifugeJSON:   centrifugeJSON,
		yafudsClient:     yafudsStorage,
		siteConfigClient: siteConfigClient,
	}
}

// todo: remove
//  Centrifuge example
func (s *Service) SendCentrifuge(ctx context.Context, centrifugeChannel string, eventData interface{}) {
	data, err := s.centrifugeJSON.Marshal(eventData)
	if err != nil {
		err = service.SemanticError{Err: err, Msg: fmt.Sprintf("unable to marshal eventData of type: %T", eventData)}.E()
		log.Error(ctx, err.Error(), "error", err)
		return
	}

	err = s.centrifugeClient.Publish(ctx, &centrifuge.PublishRequest{
		Channel: centrifugeChannel,
		Data:    data,
	})
	if err != nil {
		log.Error(ctx, err.Error(), "error", err)
	}
}
