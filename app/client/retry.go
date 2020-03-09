package client

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/context"
	gcontext "gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"

	"gitlab.com/proemergotech/log-go/v3"
)

type evalFunc func(error, *http.Request, *http.Response) error

type transport struct {
	evaluator      evalFunc
	transport      http.RoundTripper
	context        *gcontext.Context
	backoffTimeout time.Duration
	requestTimeout time.Duration
}

func RetryMiddleware(backoffTimeout time.Duration, requestTimeout time.Duration) plugin.Plugin {
	return plugin.NewPhasePlugin("before dial", func(ctx *gcontext.Context, handler gcontext.Handler) {
		intercept(ctx, backoffTimeout, requestTimeout)
		handler.Next(ctx)
	})
}

func intercept(ctx *gcontext.Context, backoffTimeout time.Duration, requestTimeout time.Duration) {
	t := &transport{
		evaluator:      evaluator,
		transport:      ctx.Client.Transport,
		context:        ctx,
		backoffTimeout: backoffTimeout,
		requestTimeout: requestTimeout,
	}
	if backoffTimeout == 0 {
		t.backoffTimeout = DefaultMaxElapsedTime
	}

	ctx.Client.Transport = t
}

var evaluator = func(err error, req *http.Request, res *http.Response) error {
	if err != nil {
		return err
	}

	if res.StatusCode >= 500 || res.StatusCode == http.StatusRequestTimeout || res.StatusCode == http.StatusNotFound {
		return errors.New("server response error")
	}

	return nil
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	res := t.context.Response

	buf, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return res, err
	}
	err = req.Body.Close()
	if err != nil {
		return res, err
	}

	b := NewExponentialBackOff(t.backoffTimeout, DefaultMaxInterval, DefaultRandomizationFactor)
	var errRetryAttempt error
	for {
		reqCopy := &http.Request{}
		*reqCopy = *req

		ctx2, cancel := context.WithTimeout(req.Context(), t.requestTimeout)
		reqCopy = reqCopy.WithContext(ctx2)

		reqCopy.Body = ioutil.NopCloser(bytes.NewBuffer(buf))

		res, err = t.transport.RoundTrip(reqCopy)
		err = t.evaluator(err, req, res)

		cancel()
		if err == nil {
			if errRetryAttempt != nil {
				log.Warn(req.Context(), "The request had to be retried.", "url", req.RequestURI, "error", errRetryAttempt)
			}
			return res, nil
		}
		errRetryAttempt = err

		hasNext, duration := b.NextBackOff()
		if !hasNext {
			return nil, err
		}

		time.Sleep(duration)
	}
}
