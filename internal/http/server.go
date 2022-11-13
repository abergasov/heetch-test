package http

import (
	"zombie_locator/internal/logger"
	"zombie_locator/internal/service/locator"

	"go.uber.org/zap"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

// Server implements a HTTP server and a router for the zombies locations endpoint.
type Server struct {
	log      logger.AppLogger
	service  locator.Locator
	appAddr  string
	fiberApp *fiber.App
}

// NewServer sets up a new Server using the provided listener address and HTTP handler for zombie locations.
func NewServer(log logger.AppLogger, address string, service locator.Locator) *Server {
	app := &Server{
		log:     log.With(zap.String("service", "http")),
		appAddr: address,
		fiberApp: fiber.New(
			fiber.Config{
				DisableStartupMessage: true,
			},
		),
		service: service,
	}
	app.fiberApp.Use(recover.New())
	app.initRoutes()
	return app
}

func (s *Server) initRoutes() {
	s.fiberApp.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.SendString("pong")
	})
	s.fiberApp.Get("/zombies", s.zombieLocationsHandler)
}

// Run starts the HTTP Server.
func (s *Server) Run() error {
	s.log.Info("Starting HTTP server on port", zap.String("port", s.appAddr))
	return s.fiberApp.Listen(s.appAddr)
}

// Shutdown gracefully shuts down the HTTP Server.
func (s *Server) Shutdown() error {
	s.log.Info("Shutting down HTTP server")
	return s.fiberApp.Shutdown()
}
