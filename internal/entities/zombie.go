package entities

import (
	"github.com/google/uuid"
)

type ZombieLocationV1 struct {
	ZombieID  uuid.UUID `json:"zombie_id"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	UpdatedAt string    `json:"updated_at"`
}

// ZombieLocationV2 future usage example
type ZombieLocationV2 struct {
	ZombieID  uuid.UUID `json:"zombie_id"`
	Type      string    `json:"type"`
	Latitude  float64   `json:"latitude"`
	Longitude float64   `json:"longitude"`
	UpdatedAt string    `json:"updated_at"`
}

type ZombieCapturedV1 struct {
	ZombieID  uuid.UUID `json:"zombie_id"`
	UpdatedAt string    `json:"updated_at"`
}
