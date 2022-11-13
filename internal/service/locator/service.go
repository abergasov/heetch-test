package locator

import (
	"context"
	"zombie_locator/internal/logger"
	"zombie_locator/internal/repository/zombie"

	"go.uber.org/zap"
)

type Service struct {
	log  logger.AppLogger
	repo zombie.Zombier
}

func NewLocatorService(log logger.AppLogger, repo zombie.Zombier) *Service {
	return &Service{
		log:  log.With(zap.String("service", "locator")),
		repo: repo,
	}
}

func (s *Service) Locate(ctx context.Context, lat, lon, limit float64) ([]zombie.Location, error) {
	return s.repo.LocateZombieList(ctx, lat, lon, limit)
}
