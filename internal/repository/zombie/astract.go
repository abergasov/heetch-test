package zombie

import (
	"context"

	"github.com/google/uuid"
)

type Location struct {
	ZombieId  uuid.UUID `json:"zombie_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
}

type Zombier interface {
	CapturedZombie(ctx context.Context, zombieId uuid.UUID, updatedAt string) error
	LocatedZombie(ctx context.Context, zombieId uuid.UUID, lat, lon float64, updatedAt string) error
	LocateZombieList(ctx context.Context, lat, lon, limitKm float64) ([]Location, error)
}
