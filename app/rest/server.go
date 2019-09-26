package rest

import (
	"context"
	"net/http"
	"time"

	"github.com/pkg/errors"
)

type Server struct {
	httpServer *http.Server
	controller *Controller
}

func NewServer(
	httpServer *http.Server,
	controller *Controller,
) *Server {
	return &Server{
		httpServer: httpServer,
		controller: controller,
	}
}

func (s *Server) Start(errorCh chan error) {
	s.controller.start()

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil {
			errorCh <- errors.Wrap(err, "http server error")
		}
	}()
}

func (s *Server) Stop(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "server graceful shutdown failed")
	}

	return nil
}
