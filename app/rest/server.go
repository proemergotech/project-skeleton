package rest

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

type Server struct {
	echoEngine *echo.Echo
	httpServer *http.Server
	router     *Router
	port       int
}

func NewServer(
	port int,
	echoEngine *echo.Echo,
	router *Router,
) *Server {
	return &Server{
		port:       port,
		echoEngine: echoEngine,
		router:     router,
	}
}

func (s *Server) Start(errorCh chan error) {
	s.router.route()

	s.httpServer = &http.Server{
		Addr:    ":" + strconv.Itoa(s.port),
		Handler: s.echoEngine,
	}

	go func() {
		err := s.httpServer.ListenAndServe()
		if err != nil {
			errorCh <- errors.Wrap(err, "http server error")
		}
	}()
}

func (s *Server) Stop(timeout time.Duration) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err = s.httpServer.Shutdown(ctx); err != nil {
		return errors.Wrap(err, "server graceful shutdown failed")
	}

	return
}
