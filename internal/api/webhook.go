package api

import "github.com/labstack/echo/v4"

// Webhook handler to receive a JSON request notifying that invoice has been paid
func (s *Server) Webhook(ctx echo.Context) error {
	return nil
}
