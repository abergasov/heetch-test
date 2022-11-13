package db

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // justifying it
)

type Connector interface {
	Client() *sqlx.DB
}
