package db

import (
	"github.com/jmoiron/sqlx"
)

type DB struct {
	conn *sqlx.DB
}

// NewRepo returns db implements Repo interface
func NewDB(conn *sqlx.DB) *DB {
	return &DB{conn: conn}
}
