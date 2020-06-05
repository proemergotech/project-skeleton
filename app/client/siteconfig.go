package client

import (
	"context"
	"net/http"
	"sync"
	"time"

	"gitlab.com/proemergotech/log-go/v3"
	"gitlab.com/proemergotech/trace-go/v2"
	"gitlab.com/proemergotech/trace-go/v2/gentlemantrace"
	"gopkg.in/h2non/gentleman.v2"

	"gitlab.com/proemergotech/dliver-project-skeleton/app/schema/siteconfig"
)

type SiteConfig struct {
	config     *siteconfig.SiteConfig
	mutex      *sync.RWMutex
	httpClient *gentleman.Client
}

func NewSiteConfig(ctx context.Context, httpClient *gentleman.Client) (*SiteConfig, error) {
	conf := &SiteConfig{httpClient: httpClient}

	resp, err := conf.loadSiteConfig(trace.WithCorrelation(ctx, trace.NewCorrelation()))
	if err != nil {
		return nil, err
	}

	conf.config = resp

	go func() {
		for range time.Tick(time.Minute) {
			resp, err := conf.loadSiteConfig(trace.WithCorrelation(context.Background(), trace.NewCorrelation()))
			if err != nil {
				log.Warn(ctx, "couldn't reload site configuration")
				continue
			}
			conf.mutex.Lock()
			conf.config = resp
			conf.mutex.Unlock()
		}
	}()

	return conf, nil
}

func (s *SiteConfig) GetSiteConfig() *siteconfig.SiteConfig {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.config
}

func (s *SiteConfig) loadSiteConfig(ctx context.Context) (*siteconfig.SiteConfig, error) {
	rawResp, err := s.request(ctx).
		Method(http.MethodGet).
		Path("/api/v1/site-config").
		Do()

	if err != nil {
		return nil, err
	}

	resp := &siteconfig.SiteConfig{}
	err = rawResp.JSON(resp)
	if err != nil {
		return nil, clientError{Err: err}.E()
	}

	return resp, nil
}

func (s *SiteConfig) request(ctx context.Context) *gentleman.Request {
	return gentlemantrace.WithContext(ctx, s.httpClient.Request())
}
