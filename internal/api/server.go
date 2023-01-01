package api

import "github.com/labstack/echo/v4"

type Server struct {
	Foo bool
}

func NewServer() *Server {
	return &Server{Foo: false}
}

// Load handler to receive a CVS document to save invoices to be paid
func (s *Server) Load(ctx echo.Context) error {
	return nil
}

// Webhook handler to receive a JSON request notifying that invoice has been paid
func (s *Server) Webhook(ctx echo.Context) error {
	return nil
}
