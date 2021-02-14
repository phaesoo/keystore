package cache

import "github.com/phaesoo/shield/pkg/memdb"

type Cache struct {
	pool *memdb.Pool
}

// NewCache returns implementation of cache layer
func NewCache(pool *memdb.Pool) *Cache {
	return &Cache{pool: pool}
}
