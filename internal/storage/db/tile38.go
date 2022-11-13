package db

import (
	"fmt"

	"github.com/xjem/t38c"
)

// NewTile38Connection creates a new DB instance based on the given configuration.
func NewTile38Connection(connStr string) (*t38c.Client, error) {
	client, err := t38c.New(connStr, t38c.Debug)
	if err != nil {
		return nil, fmt.Errorf("could not create a tile38 connection: %w", err)
	}
	if err = client.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping tile38: %w", err)
	}
	return client, nil
}
