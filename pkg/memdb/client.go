package memdb

import (
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/phaesoo/shield/configs"
	"github.com/rafaeljusto/redigomock"
)

type Client = redis.Conn

type Pool = redis.Pool

const (
	// constant db configuration
	maxIdle     int           = 256
	idleTimeout time.Duration = 120 * time.Second

	testDatabase int = 1
)

// NewPool returns redis client pool
// when a client is required, use `pool.Get()`
// and close it once it is no longer needed.
func NewPool(address string, db int) *Pool {
	return &redis.Pool{
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
		Dial: func() (Client, error) {
			client, err := newClient(address)
			if err != nil {
				return client, nil
			}
			if db != 0 {
				_, err = client.Do("SELECT", db)
			}
			return client, err
		},
	}
}

// NewMockPool creates a mock pool for testing
func NewMockPool(conn *redigomock.Conn) *Pool {
	return &redis.Pool{
		// Return the same connection mock for each Get() call.
		Dial:    func() (redis.Conn, error) { return conn, nil },
		MaxIdle: 10,
	}
}

// NewTestPool creates a connection pool for the integration test
func NewTestPool() *Pool {
	config := configs.Get().Redis
	return NewPool(config.Address(), testDatabase)
}

func newClient(address string) (Client, error) {
	return redis.Dial(
		"tcp",
		address,
		redis.DialConnectTimeout(2*time.Second),
	)
}
