package db

import (
	"fmt"

	"github.com/jmoiron/sqlx"

	// Add postgresql support
	_ "github.com/lib/pq"
)

type Connect struct {
	db *sqlx.DB
}

// NewPostgresConnection creates a new DB instance based on the given configuration.
func NewPostgresConnection(dbDSN string) (*Connect, error) {
	db, err := sqlx.Connect("postgres", dbDSN)
	if err != nil {
		return nil, fmt.Errorf("could not create a connection: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("PostgreSQL database is not reachable: %w", err)
	}
	return &Connect{db}, nil
}

func (d *Connect) Client() *sqlx.DB {
	return d.db
}
