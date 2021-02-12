package mq

import (
	"context"
	"fmt"
	"testing"

	"github.com/phaesoo/shield/pkg/memdb"
	"github.com/stretchr/testify/assert"
)

func Test_RedisMQ(t *testing.T) {
	pool := memdb.NewPool("0.0.0.0:6379", 1)
	defer pool.Close()

	rmq := NewRedisMQ(pool)

	t.Run("publish and consume", func(t *testing.T) {
		assert := assert.New(t)

		err := rmq.Publish(context.Background(), "test", []byte("123"))
		assert.NoError(err)

		val, err := rmq.Consume(context.Background(), "test")
		assert.Equal(fmt.Sprintf("%s", val), "123")
		assert.NoError(err)
	})

	t.Run("consume", func(t *testing.T) {
		assert := assert.New(t)

		val, err := rmq.Consume(context.Background(), "test")
		assert.Nil(val)
		assert.NoError(err)
	})
}
