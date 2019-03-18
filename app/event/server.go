package event

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

func (s *Server) Start() error {
	return s.controller.start()
}
