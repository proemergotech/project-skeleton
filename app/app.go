package app

import (
	"context"
	"encoding/json"

	"github.com/labstack/gommon/log"
	"github.com/pkg/errors"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/apierr"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/client/centrifugo/proto/apiproto"
)

type Core struct {
	centrifugeClient apiproto.CentrifugeClient
}

func NewCore(
	centrifugeClient apiproto.CentrifugeClient,
) *Core {
	return &Core{
		centrifugeClient: centrifugeClient,
	}
}

// Centrifuge example
func (c *Core) SendCentrifuge(ctx context.Context, namespace string, identifier string, eventData interface{}) {
	centrifugeChannel := namespace + ":" + identifier

	data, err := json.Marshal(eventData)

	if err != nil {
		err = apierr.Semantic(errors.Wrapf(err, "unable to marshal eventData of type: %T", eventData))
		log.Error(ctx, err.Error(), "error", err)
		return
	}

	resp, err := c.centrifugeClient.Publish(ctx, &apiproto.PublishRequest{
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
