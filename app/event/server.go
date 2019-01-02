package event

type Server struct {
	router *Router
}

func NewServer(
	router *Router,
) *Server {
	return &Server{
		router: router,
	}
}

func (s *Server) Start() {
	s.router.route()
}
