package client

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/pkg/errors"
	"gitlab.com/proemergotech/dliver-project-skeleton/errorsf"
	"gitlab.com/proemergotech/log-go/v3"
	"golang.org/x/net/context"
	gcontext "gopkg.in/h2non/gentleman.v2/context"
	"gopkg.in/h2non/gentleman.v2/plugin"
)

var noopCancel = func() {}

type Option func(*transport)

type evalFunc func(error, *http.Request, *http.Response) (retry bool, err error)

type transport struct {
	evaluator      evalFunc
	transport      http.RoundTripper
	gctx           *gcontext.Context
	backoffTimeout time.Duration
	requestTimeout time.Duration
	loggingEnabled bool
	logResponse    bool
}

func RetryMiddleware(options ...Option) plugin.Plugin {
	return plugin.NewPhasePlugin("before dial", func(gctx *gcontext.Context, handler gcontext.Handler) {
		t := &transport{
			evaluator:      defaultEvaluator,
			transport:      gctx.Client.Transport,
			gctx:           gctx,
			backoffTimeout: DefaultMaxElapsedTime,
		}

		for _, opt := range options {
			opt(t)
		}

		gctx.Client.Transport = t

		handler.Next(gctx)
	})
}

func (t *transport) RoundTrip(req *http.Request) (*http.Response, error) {
	var body []byte

	if req.Body != nil {
		buf, err := ioutil.ReadAll(req.Body)
		if err != nil {
			_ = req.Body.Close()
			return t.gctx.Response, err
		}
		_ = req.Body.Close()

		body = buf
	}

	reqCopy := req
	resetTimeout := false

	if t.requestTimeout > 0 {
		if _, ok := req.Context().Deadline(); !ok {
			reqCopy = req.Clone(req.Context())
			resetTimeout = true
		}
	}

	return t.retry(reqCopy, body, reqCopy.Context().Done(), resetTimeout)
}

func (t *transport) retry(req *http.Request, body []byte, done <-chan struct{}, resetTimeout bool) (*http.Response, error) {
	cancel := noopCancel
	retryCount := 0
	backoff := NewExponentialBackOff(t.backoffTimeout, DefaultMaxInterval, DefaultRandomizationFactor)
	origCtx := req.Context()

	for {
		if resetTimeout {
			ctx, c := context.WithTimeout(origCtx, t.requestTimeout)
			req = req.WithContext(ctx)
			cancel = c
		}

		if body != nil {
			req.Body = ioutil.NopCloser(bytes.NewBuffer(body))
		}

		res, err := t.transport.RoundTrip(req)
		retry, err := t.evaluator(err, req, res)
		if !retry {
			return res, err
		}

		cancel()

		if res != nil {
			if t.logResponse {
				b, _ := httputil.DumpResponse(res, true)
				_ = res.Body.Close()

				err = errorsf.WithFields(err, "failed_retry_response", string(b))
			} else {
				_, _ = io.Copy(ioutil.Discard, res.Body)
				_ = res.Body.Close()
			}
		}

		hasNext, duration := backoff.NextBackOff()
		if !hasNext {
			return nil, err
		}

		select {
		case <-time.After(duration):
			retryCount++
			if t.loggingEnabled {
				log.Warn(req.Context(), fmt.Sprintf("error during request, retry # %d", retryCount), "error", err)
			}
		case <-done:
			return nil, err
		}
	}
}

func defaultEvaluator(err error, req *http.Request, res *http.Response) (bool, error) {
	if err != nil {
		return true, err
	}

	if res.StatusCode >= 500 || res.StatusCode == http.StatusRequestTimeout {
		return true, errors.New("server response error")
	}

	return false, nil
}

func BackoffTimeout(timeout time.Duration) Option {
	return func(t *transport) {
		t.backoffTimeout = timeout
	}
}

func RequestTimeout(timeout time.Duration) Option {
	return func(t *transport) {
		t.requestTimeout = timeout
	}
}

func Evaluator(evalFn evalFunc) Option {
	return func(t *transport) {
		t.evaluator = evalFn
	}
}

func EnableLogging() Option {
	return func(t *transport) {
		t.loggingEnabled = true
	}
}
func LogResponse() Option {
	return func(t *transport) {
		t.logResponse = true
	}
}
