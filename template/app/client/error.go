package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/proemergotech/errors"
	"github.com/proemergotech/log/v3"
	gcontext "gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"

	//%:{{ `
	"github.com/proemergotech/project-skeleton/app/schema"
	"github.com/proemergotech/project-skeleton/app/schema/skeleton"
	//%: ` | replace "project-skeleton" .ProjectName | replace "skeleton" .SchemaPackage }}
)

const (
	errCode     = "remote_code"
	errDetails  = "remote_details"
	errHTTPCode = "remote_http_code"
	errService  = "remote_service"
)

func RestErrorMiddleware(serviceName string) plugin.Plugin {
	afterDialHandler := func(gCtx *gcontext.Context, handler gcontext.Handler) {
		res := gCtx.Response

		if res == nil || res.StatusCode < 400 {
			handler.Next(gCtx)
			return
		}

		defer func() {
			if err := res.Body.Close(); err != nil {
				//%:{{ `
				err = skeleton.SemanticError{Err: err}.E()
				//%: ` | replace "skeleton" .SchemaPackage | trim }}
				log.Error(gCtx.Request.Context(), err.Error(), "error", err)
			}
		}()

		handler.Error(gCtx, clientHTTPError{Res: res, ServiceName: serviceName}.E())
	}

	errorHandler := func(gCtx *gcontext.Context, handler gcontext.Handler) {
		if gCtx.Error == nil {
			handler.Next(gCtx)
			return
		}

		err := gCtx.Error
		if schema.ErrorCode(err) == "" {
			err = clientError{Err: err, Msg: fmt.Sprintf("error occured during call of %s", serviceName)}.E()
		}

		handler.Error(gCtx, err)
	}

	handlers := plugin.Handlers{
		"after dial": afterDialHandler,
		"error":      errorHandler,
	}

	return &plugin.Layer{
		Handlers: handlers,
	}
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
		return errors.WithFields(
			errors.New("invalid error response format"),
			errHTTPCode, e.Res.StatusCode,
			errService, e.ServiceName,
			//%:{{ `
			schema.ErrCode, skeleton.ErrClient,
			//%: ` | replace "skeleton" .SchemaPackage | trim }}
			schema.ErrHTTPCode, 500,
		)
	} else if errResp.Error.Code == "" {
		return errors.WithFields(
			errors.New("invalid error response format: missing code field"),
			errHTTPCode, e.Res.StatusCode,
			errService, e.ServiceName,
			//%:{{ `
			schema.ErrCode, skeleton.ErrClient,
			//%: ` | replace "skeleton" .SchemaPackage | trim }}
			schema.ErrHTTPCode, 500,
		)
	}

	msg := fmt.Sprintf("error calling %s service", e.ServiceName)
	if errResp.Error.Message != "" {
		msg = fmt.Sprintf("%s: %s", msg, errResp.Error.Message)
	}

	return errors.WithFields(
		errors.New(msg),
		errCode, errResp.Error.Code,
		errHTTPCode, e.Res.StatusCode,
		errDetails, errResp.Error.Details,
		errService, e.ServiceName,
		//%:{{ `
		schema.ErrCode, skeleton.ErrClient,
		//%: ` | replace "skeleton" .SchemaPackage | trim }}
		schema.ErrHTTPCode, 500,
	)
}

type clientError struct {
	Err error
	Msg string
}

func (e clientError) E() error {
	msg := "client error"

	if e.Msg != "" {
		msg = e.Msg
	}

	return errors.WithFields(
		errors.WrapOrNew(e.Err, msg),
		//%:{{ `
		schema.ErrCode, skeleton.ErrClient,
		//%: ` | replace "skeleton" .SchemaPackage | trim }}
		schema.ErrHTTPCode, 500,
	)
}

func ErrorCode(err error) string {
	field := errors.Field(err, errCode)
	if field == nil {
		return ""
	}

	return field.(string)
}

func ErrorHTTPCode(err error) int {
	field := errors.Field(err, errHTTPCode)
	if field == nil {
		return 0
	}

	return field.(int)
}

func ErrorDetails(err error) []map[string]interface{} {
	field := errors.Field(err, errDetails)
	if field == nil {
		return []map[string]interface{}{}
	}

	return field.([]map[string]interface{})
}
