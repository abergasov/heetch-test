package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type zombieLocationPayload struct {
	Lat   float64 `json:"lat"`
	Lon   float64 `json:"lon"`
	Limit float64 `json:"limit"`
}

var wrongPayload = errors.New("lat and lon must be greater than 0")

// zombieLocationsHandler processes HTTP requests for zombie locations.
func (s *Server) zombieLocationsHandler(ctx *fiber.Ctx) error {
	log := s.log.
		With(zap.String("method", "zombieLocationsHandler")).
		With(zap.ByteString("query", ctx.Request().URI().QueryString()))
	var payload zombieLocationPayload
	if err := ctx.QueryParser(&payload); err != nil {
		log.Error("failed to parse query", err)
		return fiber.ErrBadRequest
	}

	if payload.Lat <= 0 || payload.Lon <= 0 {
		log.Error("invalid query parameters", wrongPayload)
		return fiber.ErrBadRequest
	}

	data, err := s.service.Locate(ctx.UserContext(), payload.Lat, payload.Lon, payload.Limit)
	if err != nil {
		log.Error("failed to locate zombies", err)
		return fiber.ErrInternalServerError
	}
	return ctx.JSON(data)
}
