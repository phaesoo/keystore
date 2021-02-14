package repo

import (
	"github.com/jmoiron/sqlx"
)

type db struct {
	conn *sqlx.DB
}

type Repo interface {
	keyauthRepo
}

// NewRepo returns db implements Repo interface
func NewRepo(conn *sqlx.DB) Repo {
	return &db{conn: conn}
}
