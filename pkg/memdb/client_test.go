package memdb

import (
	"fmt"
	"testing"

	"github.com/gomodule/redigo/redis"
	"github.com/phaesoo/shield/configs"
	"github.com/stretchr/testify/assert"
)

func Test_NewPool(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	conf := configs.Get().Redis

	t.Run("it dials the configured host", func(t *testing.T) {
		pool := NewPool(conf.Address(), conf.Database)
		defer pool.Close()

		c := pool.Get()
		err := c.Send("PING")
		assert.NoError(t, err)
		c.Close()
	})

	t.Run("it uses the correct database", func(t *testing.T) {
		conf := conf
		conf.Database = testDatabase

		assert := assert.New(t)
		c, err := newClient(conf.Address())
		assert.NoError(err)

		_, err = c.Do("SET", "hello", "world")
		assert.NoError(err)
		_, err = c.Do("SELECT", conf.Database)
		assert.NoError(err)
		_, err = c.Do("SET", "hello", "bye")
		assert.NoError(err)
		c.Close()

		pool := NewPool(conf.Address(), conf.Database)
		defer pool.Close()

		c2 := pool.Get()
		defer c2.Close()

		v, err := redis.String(c2.Do("GET", "hello"))
		assert.NoError(err)
		fmt.Println(v, err)
		assert.Equal("bye", v)
	})
}
