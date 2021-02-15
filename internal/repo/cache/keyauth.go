package cache

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"github.com/phaesoo/shield/internal/models"
	"github.com/pkg/errors"
)

const (
	authKeyHashPrefix  = "auth-key:"
	pathPermissionHash = "path-permissions"
)

func (c *Cache) AuthKey(ctx context.Context, accessKey string) (models.AuthKey, error) {
	conn := c.pool.Get()
	defer conn.Close()

	val, err := redis.Values(conn.Do("HGETALL", authKeyHash(accessKey)))
	if err != nil {
		return models.AuthKey{}, errors.Wrap(err, "Get auth key from redis")
	}

	var authKey models.AuthKey
	if err := redis.ScanStruct(val, &authKey); err != nil {
		return authKey, err
	} else if authKey == (models.AuthKey{}) {
		return authKey, ErrNotFound
	}

	return authKey, nil
}

func (c *Cache) SetAuthKey(ctx context.Context, authKey models.AuthKey) error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("HMSET", redis.Args{authKeyHash(authKey.AccessKey)}.AddFlat(authKey)...)
	if err != nil {
		return err
	}
	return nil
}

func authKeyHash(accessKey string) string {
	return authKeyHashPrefix + accessKey
}

func (c *Cache) PathPermission(ctx context.Context, id int) (models.PathPermission, error) {
	conn := c.pool.Get()
	defer conn.Close()

	pathPattern, err := redis.String(conn.Do("HGET", pathPermissionHash, id))
	if err == redis.ErrNil {
		return models.PathPermission{}, ErrNotFound
	} else if err != nil {
		return models.PathPermission{}, errors.Wrap(err, "Get path permission")
	}

	return models.PathPermission{ID: id, PathPattern: pathPattern}, nil
}

func (c *Cache) RefreshPathPermissions(ctx context.Context, perms []models.PathPermission) error {
	conn := c.pool.Get()
	defer conn.Close()

	_ = conn.Send("MULTI")
	_ = conn.Send("DEL", pathPermissionHash)

	for _, p := range perms {
		_ = conn.Send("HSET", pathPermissionHash, p.ID, p.PathPattern)
	}

	if _, err := conn.Do("EXEC"); err != nil {
		return errors.Wrap(err, "Refresh path permissions")
	}
	return nil
}

func (c *Cache) PathPermissionIDs(ctx context.Context, accessKey string) ([]int, error) {
	conn := c.pool.Get()
	defer conn.Close()

	var permIDs []int
	permIDs, err := redis.Ints(conn.Do("LRANGE", accessKey, 0, -1))
	if err == redis.ErrNil {
		return permIDs, ErrNotFound
	} else if err != nil {
		return permIDs, errors.Wrap(err, "Get path permission IDs")
	}
	return permIDs, nil
}

func (c *Cache) SetPathPermissionIDs(ctx context.Context, accessKey string, permIDs []int) error {
	conn := c.pool.Get()
	defer conn.Close()

	_, err := conn.Do("LPUSH", accessKey, permIDs)
	if err != nil {
		return errors.Wrap(err, "Set path permission IDs")
	}
	return nil
}
