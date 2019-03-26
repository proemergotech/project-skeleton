package client

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/pkg/errors"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema"
	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/service"
	"gitlab.com/proemergotech/dliver-project-skeleton/errorsf"
	"gitlab.com/proemergotech/log-go"
	gcontext "gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
)

const (
	errCode     = "remote_code"
	errDetails  = "remote_details"
	errHTTPCode = "remote_http_code"
	errService  = "remote_service"
)

func restErrorMiddleware(serviceName string) plugin.Plugin {
	return plugin.NewPhasePlugin("after dial", func(gCtx *gcontext.Context, handler gcontext.Handler) {
		res := gCtx.Response

		if res == nil || res.StatusCode < 400 {
			handler.Next(gCtx)
			return
		}

		defer func() {
			err := res.Body.Close()
			if err != nil {
				err = service.SemanticError{Err: err}.E()
				log.Error(gCtx.Request.Context(), err.Error(), "error", err)
			}
		}()

		handler.Error(gCtx, clientHTTPError{Res: res, ServiceName: serviceName}.E())
	})
}

type clientHTTPError struct {
	Res         *http.Response
	ServiceName string
}

func (e clientHTTPError) E() error {
	jsonDecoder := json.NewDecoder(e.Res.Body)

	errResp := &schema.HTTPError{}
	err := jsonDecoder.Decode(errResp)
	if err != nil && err != io.EOF {
		return errorsf.WithFields(
			errors.New("invalid error response format"),
			errHTTPCode, e.Res.StatusCode,
			errService, e.ServiceName,
		)
	} else if errResp.Error.Code == "" {
		return errorsf.WithFields(
			errors.New("invalid error response format: missing code field"),
			errHTTPCode, e.Res.StatusCode,
			errService, e.ServiceName,
		)
	}

	return errorsf.WithFields(
		errors.New(errResp.Error.Message),
		errCode, errResp.Error.Code,
		errHTTPCode, e.Res.StatusCode,
		errDetails, errResp.Error.Details,
		errService, e.ServiceName,
	)
}

func ErrorCode(err error) string {
	field := errorsf.Field(err, errCode)
	if field == nil {
		return ""
	}

	return field.(string)
}

func ErrorHTTPCode(err error) int {
	field := errorsf.Field(err, errHTTPCode)
	if field == nil {
		return 0
	}

	return field.(int)
}

func ErrorDetails(err error) []map[string]interface{} {
	field := errorsf.Field(err, errDetails)
	if field == nil {
		return []map[string]interface{}{}
	}

	return field.([]map[string]interface{})
}
