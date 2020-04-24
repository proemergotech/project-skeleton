package client

import (
	"context"
	"net/http"

	"gitlab.com/proemergotech/trace-go/v2/gentlemantrace"
	"gopkg.in/h2non/gentleman.v2"
)

// todo: remove
//  Example client
type Dummy struct {
	httpClient *gentleman.Client
}

func NewDummy(httpClient *gentleman.Client) *Dummy {
	return &Dummy{
		httpClient: httpClient,
	}
}

func (d *Dummy) GetDummy(ctx context.Context) error {
	resp, err := d.request(ctx).
		Method(http.MethodGet).
		Path("").
		Do()
	if err != nil {
		return err
	}

	res := &struct{}{}
	if err := resp.JSON(res); err != nil {
		return clientError{Err: err, Msg: "dummy service error: failed to bind response"}.E()
	}

	return nil
}

func (d *Dummy) request(ctx context.Context) *gentleman.Request {
	return gentlemantrace.WithContext(ctx, d.httpClient.Request())
}
