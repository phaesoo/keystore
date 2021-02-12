package store

import "github.com/phaesoo/shield/pkg/memdb"

type db struct {
	pool *memdb.Pool
}

type Store interface {
	commonStore
	masterStore
}

// NewStore returns db implements Store interface
func NewStore(pool *memdb.Pool) Store {
	return &db{pool: pool}
}
