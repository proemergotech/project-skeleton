package client

import (
	"encoding/json"
	"io"

	"github.com/pkg/errors"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/apierr"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/log-go"
	gcontext "gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
)

func restErrorMiddleware(errorPrefix string) plugin.Plugin {
	return plugin.NewPhasePlugin("after dial", func(gCtx *gcontext.Context, handler gcontext.Handler) {
		resp := gCtx.Response

		if resp == nil || resp.StatusCode < 400 {
			handler.Next(gCtx)
			return
		}

		jsonDecoder := json.NewDecoder(resp.Body)
		defer func() {
			err := resp.Body.Close()
			if err != nil {
				err = apierr.Semantic(err)
				log.Error(gCtx.Request.Context(), err.Error(), "error", err)
			}
		}()

		errResp := &schema.HTTPError{Error: schema.NewError(errorPrefix)}
		err := jsonDecoder.Decode(errResp)
		if err != nil && err != io.EOF {
			handler.Error(gCtx, errors.Wrap(err, "Invalid error response format"))
			return
		}
		if errResp.Error.Code == "" {
			handler.Error(gCtx, errors.New("Invalid error response format: missing code field"))
			return
		}

		handler.Error(gCtx, errResp.Error)
	})
}
