//%: {{ if .Geb }}
package event

import (
	"time"
)

type Server struct {
	controller *Controller
}

func NewServer(
	controller *Controller,
) *Server {
	return &Server{
		controller: controller,
	}
}

func (s *Server) Start(errorCh chan<- error) {
	err := s.controller.start()
	if err != nil {
		go func() {
			errorCh <- err
		}()
	}
}

func (s *Server) Stop(_ time.Duration) error {
	return nil
}

//%: {{ end }}
