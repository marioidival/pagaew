package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/marioidival/pagaew/internal/api"
)

type Server interface {
	Webhook(ctx echo.Context) error
	Load(ctx echo.Context) error
}

// Setup create the basic handlers to API
func Setup() *echo.Echo {
	e := echo.New()
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(500)))

	server := api.NewServer()

	registerHandlers(e, server)

	return e
}


func registerHandlers(router *echo.Echo, server Server) {
	router.POST("/XXXX", server.Webhook)
	router.POST("/YYYY", server.Load)
}
