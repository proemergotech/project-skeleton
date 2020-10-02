package service

import (
	"context"
	"fmt"

	jsoniter "github.com/json-iterator/go"
	"gitlab.com/proemergotech/log-go/v3"

	//%:{{ `
	"gitlab.com/proemergotech/dliver-project-skeleton/app/client"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/centrifuge"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/skeleton"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/storage"
	//%: ` | replace "dliver-project-skeleton" .ProjectName | replace "skeleton" .SchemaPackage }}
)

type Service struct {
	//%: {{- if .Centrifuge }}
	centrifugeClient *client.Centrifuge
	centrifugeJSON   jsoniter.API
	//%: {{- end }}
	//%: {{- if .Yafuds }}
	yafudsClient *storage.Yafuds
	//%: {{- end }}
	//%: {{- if .SiteConfig }}
	siteConfigClient *client.SiteConfig
	//%: {{- end }}
}

func NewService(
	//%: {{- if .Centrifuge }}
	centrifugeClient *client.Centrifuge,
	centrifugeJSON jsoniter.API,
	//%: {{- end }}
	//%: {{- if .Yafuds }}
	yafudsStorage *storage.Yafuds,
	//%: {{- end }}
	//%: {{- if .SiteConfig }}
	siteConfigClient *client.SiteConfig,
	//%: {{- end }}
) *Service {
	return &Service{
		//%: {{- if .Centrifuge }}
		centrifugeClient: centrifugeClient,
		centrifugeJSON:   centrifugeJSON,
		//%: {{- end }}
		//%: {{- if .Yafuds }}
		yafudsClient: yafudsStorage,
		//%: {{- end }}
		//%: {{- if .SiteConfig }}
		siteConfigClient: siteConfigClient,
		//%: {{- end }}
	}
}

//%: {{ if .Examples }}{{ `
func (s *Service) Dummy(ctx context.Context, req *skeleton.DummyRequest) error {
	//%: ` | replace "skeleton" .SchemaPackage }}
	//%: {{- if .Centrifuge }}{{ `
	s.SendCentrifuge(ctx, centrifuge.GlobalChannel("dummy", "dummy_site_group", "dummy"), skeleton.CentrifugeData{
		Group: req.DummyGroup,
		UUID:  req.DummyUUID,
	})
	//%: ` | replace "skeleton" .SchemaPackage }}{{- end }}

	return nil
} //%: {{ end }}

//%: {{- if and .Centrifuge .Examples }}
// todo: remove
//  Centrifuge example
func (s *Service) SendCentrifuge(ctx context.Context, centrifugeChannel string, eventData interface{}) {
	data, err := s.centrifugeJSON.Marshal(eventData)
	if err != nil {
		//%:{{ `
		err = skeleton.SemanticError{Err: err, Msg: fmt.Sprintf("unable to marshal eventData of type: %T", eventData)}.E()
		//%: ` | replace "skeleton" .SchemaPackage }}
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

//%: {{- end }}
