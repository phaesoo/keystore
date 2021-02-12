package store

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"github.com/phaesoo/shield/pkg/memdb"
	"github.com/pkg/errors"
)

type commonStore interface {
	Set(ctx context.Context, key string, val interface{}) error
	SetWithExpiration(ctx context.Context, key string, val interface{}, ttl uint64) error
	GetInt64(ctx context.Context, key string) (int64, error)
	GetString(ctx context.Context, key string) (string, error)
	IncrBy(ctx context.Context, key string, incr int64) (int64, error)
	GetAutoIncrID(ctx context.Context, key string) (int64, error)
	SetNX(ctx context.Context, key string, val interface{}) error
	Exists(ctx context.Context, key string) (bool, error)
	AddToSortedSet(ctx context.Context, key string, score int64, member interface{}) (int64, error)
	QueryFromSortedSet(ctx context.Context, key string, offset int, limit int) ([]interface{}, error)
	PopAndPushSortedSet(ctx context.Context, keySrc string, keyDest string, size int) ([]interface{}, error)
	RemoveFromSortedSet(ctx context.Context, key string, member interface{}) (int64, error)
	PopAndPushSortedSetWithMember(ctx context.Context, keySrc string, keyDest string, member interface{}, newScore int64) (bool, error)
	MoveSetMembersWithinScore(ctx context.Context, keySrc string, keyDest string, minScore int64, maxScore int64, newScoreIncr int64) (int64, error)
}

func (db *db) Set(ctx context.Context, key string, value interface{}) error {
	c := db.pool.Get()
	defer c.Close()

	_, err := c.Do("SET", key, value)
	if err != nil {
		return errors.Wrapf(err, "failed redis command SET %s %v", key, value)
	}
	return nil
}

func (db *db) SetWithExpiration(ctx context.Context, key string, value interface{}, ttl uint64) error {
	c := db.pool.Get()
	defer c.Close()

	_, err := c.Do("SETEX", key, ttl, value)
	if err != nil {
		return errors.Wrapf(err, "failed redis command SETEX %s %v %d", key, value, ttl)
	}
	return nil
}

func (db *db) GetInt64(ctx context.Context, key string) (int64, error) {
	c := db.pool.Get()
	defer c.Close()

	val, err := redis.Int64(c.Do("GET", key))
	if err != nil {
		return 0, errors.Wrapf(err, "failed redis command Get(int64) %s", key)
	}
	return val, nil
}

func (db *db) GetString(ctx context.Context, key string) (string, error) {
	c := db.pool.Get()
	defer c.Close()

	val, err := redis.String(c.Do("GET", key))
	if err != nil {
		return "", errors.Wrapf(err, "failed redis command Get(string) %s", key)
	}
	return val, nil
}

func (db *db) IncrBy(ctx context.Context, key string, value int64) (int64, error) {
	c := db.pool.Get()
	defer c.Close()

	val, err := redis.Int64(c.Do("INCRBY", key, value))
	if err != nil {
		return 0, errors.Wrapf(err, "failed redis command INCRBY %s %v", key, value)
	}
	return val, nil
}

func (db *db) GetAutoIncrID(ctx context.Context, key string) (int64, error) {
	return db.IncrBy(ctx, key, 1)
}

func (db *db) SetNX(ctx context.Context, key string, value interface{}) error {
	c := db.pool.Get()
	defer c.Close()

	_, err := c.Do("SETNX", key, value)
	if err != nil {
		return errors.Wrapf(err, "failed redis command SETNX %s %v", key, value)
	}
	return nil
}

func (db *db) Exists(ctx context.Context, key string) (bool, error) {
	c := db.pool.Get()
	defer c.Close()

	val, err := redis.Bool(c.Do("EXISTS", key))
	if err != nil {
		return false, errors.Wrapf(err, "failed redis command Exists(string) %s", key)
	}
	return val, nil
}

func (db *db) AddToSortedSet(ctx context.Context, key string, score int64, member interface{}) (int64, error) {
	c := db.pool.Get()
	defer c.Close()

	val, err := redis.Int64(c.Do("ZADD", key, score, member))
	if err != nil {
		return 0, errors.Wrapf(err, "failed redis command ZADD %s XX %v %v", key, score, member)
	}
	return val, nil
}

func (db *db) QueryFromSortedSet(ctx context.Context, key string, offset int, limit int) ([]interface{}, error) {
	c := db.pool.Get()
	defer c.Close()

	val, err := redis.Values(c.Do("ZRANGE", key, offset, limit))
	if err != nil {
		return nil, errors.Wrapf(err, "failed redis command ZRANGE %s %d %d", key, offset, limit)
	}
	return val, nil
}

func (db *db) PopAndPushSortedSet(ctx context.Context, keySrc string, keyDest string, size int) ([]interface{}, error) {
	c := db.pool.Get()
	defer c.Close()

	val, err := redis.Values(memdb.PopAndPushSortedSet.Do(c, keySrc, keyDest, size))
	if err != nil {
		return nil, errors.Wrapf(err, "failed redis command PopAndPushSortedSet %s %s %d", keySrc, keyDest, size)
	}
	return val, nil
}

func (db *db) PopAndPushSortedSetWithMember(ctx context.Context, keySrc string, keyDest string, member interface{}, newScore int64) (bool, error) {
	c := db.pool.Get()
	defer c.Close()

	processed, err := redis.Bool(memdb.PopAndPushSortedSetWithMember.Do(c, keySrc, keyDest, member, newScore))
	if err != nil {
		return false, errors.Wrapf(err, "failed redis command PopAndPushSortedSetWithMember %s %s %v %d", keySrc, keyDest, member, newScore)
	}
	return processed, nil
}
func (db *db) MoveSetMembersWithinScore(ctx context.Context, keySrc string, keyDest string, minScore int64, maxScore int64, newScoreIncr int64) (int64, error) {
	c := db.pool.Get()
	defer c.Close()

	numMoved, err := redis.Int64(memdb.MoveSetMembersWithinScore.Do(c, keySrc, keyDest, minScore, maxScore, newScoreIncr))
	if err != nil {
		return 0, errors.Wrapf(err, "failed redis command MoveSetMembersWithinScore %s %s %d %d %d", keySrc, keyDest, minScore, maxScore, newScoreIncr)
	}
	return numMoved, nil
}

func (db *db) RemoveFromSortedSet(ctx context.Context, key string, member interface{}) (int64, error) {
	c := db.pool.Get()
	defer c.Close()

	val, err := redis.Int64(c.Do("ZREM", key, member))
	if err != nil {
		return 0, errors.Wrapf(err, "failed redis command ZREM %s %v", key, member)
	}
	return val, nil
}
