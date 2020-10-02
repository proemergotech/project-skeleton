//%: {{ if .Centrifuge }}
package client

import (
	"context"
	"net/http"

	"gitlab.com/proemergotech/bind/gentlemanbind"
	"gitlab.com/proemergotech/trace-go/v2/gentlemantrace"
	"gopkg.in/h2non/gentleman.v2"

	//%:{{ `
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/centrifuge"
	//%: ` | replace "dliver-project-skeleton" .ProjectName }}
)

type Centrifuge struct {
	httpClient *gentleman.Client
}

func NewCentrifuge(httpClient *gentleman.Client) *Centrifuge {
	return &Centrifuge{httpClient: httpClient}
}

func (c *Centrifuge) Publish(ctx context.Context, req *centrifuge.PublishRequest) error {
	_, err := c.request(ctx).
		Method(http.MethodPost).
		Path("/api/v1/publish").
		Use(gentlemanbind.Bind(req)).
		Do()
	if err != nil {
		return err
	}

	return nil
}

func (c *Centrifuge) request(ctx context.Context) *gentleman.Request {
	return gentlemantrace.WithContext(ctx, c.httpClient.Request())
}

//%: {{ end }}
