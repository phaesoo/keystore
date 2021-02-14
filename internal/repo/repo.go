package repo

import (
	"github.com/jmoiron/sqlx"
	"github.com/phaesoo/shield/internal/repo/cache"
	"github.com/phaesoo/shield/internal/repo/db"
	"github.com/phaesoo/shield/pkg/memdb"
)

type repo struct {
	db    *db.DB
	cache *cache.Cache
}

type Repo interface {
	keyauthRepo
}

// NewRepo returns db implements Repo interface
func NewRepo(conn *sqlx.DB, pool *memdb.Pool) Repo {
	return &repo{
		db:    db.NewDB(conn),
		cache: cache.NewCache(pool),
	}
}
