package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/marioidival/pagaew/internal/api"
	"github.com/marioidival/pagaew/internal/repository"
	"github.com/marioidival/pagaew/pkg/database"
)

type Server interface {
	Webhook(ctx echo.Context) error
	Load(ctx echo.Context) error
}

// Setup create the basic handlers to API
func Setup(dbc *database.Client, productionEnv bool) *echo.Echo {
	e := echo.New()
	e.Use(middleware.RateLimiter(middleware.NewRateLimiterMemoryStore(500)))

	i := repository.NewInvoiceMemoryRepository()
	l := repository.NewLogMemoryRepository()

	server := api.NewServer(i, l)

	registerHandlers(e, server)

	return e
}


func registerHandlers(router *echo.Echo, server Server) {
	router.POST("/webhook", server.Webhook)
	router.POST("/load", server.Load)
}
