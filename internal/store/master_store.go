package store

import (
	"context"
	"encoding/json"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

type masterStore interface {
	Master(ctx context.Context) (Master, error)
	UpdateMaster(ctx context.Context, master Master) error
}

const (
	masterKey = "master"
)

func (db *db) Master(ctx context.Context) (Master, error) {
	c := db.pool.Get()
	defer c.Close()

	b, err := redis.Bytes(c.Do("GET", masterKey))
	if err != nil {
		return Master{}, errors.Wrap(err, "could not retrieve Master")
	}

	var master Master
	json.Unmarshal(b, &master)
	return master, errors.Wrap(err, "could not unmarshal Master")
}

func (db *db) UpdateMaster(ctx context.Context, master Master) error {
	c := db.pool.Get()
	defer c.Close()

	b, err := json.Marshal(master)
	if err != nil {
		return errors.Wrap(err, "could not marshal market")
	}

	_, err = c.Do("SET", masterKey, b)
	return errors.Wrap(err, "could not set Master")
}
