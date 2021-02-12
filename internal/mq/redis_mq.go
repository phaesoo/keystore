package mq

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"github.com/phaesoo/shield/pkg/memdb"
)

// RedisMQ is Redis implementation of message queue
type RedisMQ struct {
	pool *memdb.Pool
}

func NewRedisMQ(pool *memdb.Pool) *RedisMQ {
	return &RedisMQ{pool: pool}
}

func (mq *RedisMQ) Publish(ctx context.Context, topic string, data []byte) error {
	c := mq.pool.Get()
	defer c.Close()

	_, err := c.Do("LPUSH", topic, data)
	if err != nil {
		return err
	}
	return nil
}

func (mq *RedisMQ) Consume(ctx context.Context, topic string) ([]byte, error) {
	c := mq.pool.Get()
	defer c.Close()

	data, err := redis.Bytes(c.Do("RPOP", topic))
	if err != nil {
		if err == redis.ErrNil {
			return nil, nil
		}
		return nil, err
	}

	return data, nil
}
