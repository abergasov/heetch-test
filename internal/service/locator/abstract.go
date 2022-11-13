package locator

import (
	"context"
	"zombie_locator/internal/repository/zombie"
)

//go:generate mockgen -source=abstract.go -destination=abstract_locator_mock.go -package=locator
type Locator interface {
	// Locate returns the location of the zombies near provided coordinates.
	Locate(ctx context.Context, lat, lon, limit float64) ([]zombie.Location, error)
}
