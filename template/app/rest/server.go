package rest

import (
	"context"
	"time"

	"github.com/labstack/echo/v4"
	"gitlab.com/proemergotech/errors"
)

type Controller interface {
	Start()
}

type Server struct {
	echoEngine *echo.Echo
	controller Controller
}

func NewServer(
	echoEngine *echo.Echo,
	controller Controller,
) *Server {
	return &Server{
		echoEngine: echoEngine,
		controller: controller,
	}
}

func (s *Server) Start(errorCh chan<- error) {
	s.controller.Start()

	go func() {
		if err := s.echoEngine.StartServer(s.echoEngine.Server); err != nil {
			errorCh <- errors.Wrap(err, "http server error")
		}
	}()
}

func (s *Server) Stop(timeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := s.echoEngine.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "server graceful shutdown failed")
	}

	return nil
}
