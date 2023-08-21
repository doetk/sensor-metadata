package server

import "github.com/gofiber/fiber/v2"

type Server struct {
	app *fiber.App
}

func NewServer(cfg ...fiber.Config) *Server {
	return &Server{
		app: fiber.New(cfg...),
	}
}

// Listen serves HTTP requests on the given addr
//
//	s.Listen(":8080")
//	s.Listen("127.0.0.1:8080")
func (s *Server) Listen(addr string) error {
	return s.app.Listen(addr)
}

// Shutdown works by first closing all open listeners, waits indefinitely for all connections to return to idle and then shut down.
//
// Make sure the program doesn't exit and waits instead for Shutdown to return.
//
// Shutdown does not close keepalive connections, so It's recommended to set ReadTimeout to something else than 0.
func (s *Server) Shutdown() error {
	return s.app.Shutdown()
}
