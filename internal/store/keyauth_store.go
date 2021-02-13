package store

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

type keyauthStore interface {
	AuthKey(ctx context.Context, accessKey string) (models.AuthKey, error)
	SetAuthKey(ctx context.Context, authKey models.AuthKey) error
	PathPermission(ctx context.Context, id int) (models.PathPermission, error)
	RefreshPathPermissions(ctx context.Context, perms []models.PathPermission) error
}

func (db *db) AuthKey(ctx context.Context, accessKey string) (models.AuthKey, error) {
	c := db.pool.Get()
	defer c.Close()

	val, err := redis.Values(c.Do("HGETALL", authKeyHash(accessKey)))
	if err != nil {
		return models.AuthKey{}, errors.Wrap(err, "could not get ticker from redis")
	}

	var authKey models.AuthKey
	if err := redis.ScanStruct(val, &authKey); err != nil {
		return authKey, err
	} else if authKey == (models.AuthKey{}) {
		return authKey, ErrNotFound
	}

	return authKey, nil
}

func (db *db) SetAuthKey(ctx context.Context, authKey models.AuthKey) error {
	c := db.pool.Get()
	defer c.Close()

	_, err := c.Do("HMSET", redis.Args{authKeyHash(authKey.AccessKey)}.AddFlat(authKey)...)
	if err != nil {
		return err
	}
	return nil
}

func authKeyHash(accessKey string) string {
	return authKeyHashPrefix + accessKey
}

func (db *db) PathPermission(ctx context.Context, id int) (models.PathPermission, error) {
	c := db.pool.Get()
	defer c.Close()

	pathPattern, err := redis.String(c.Do("HGET", pathPermissionHash, id))
	if err == redis.ErrNil {
		return models.PathPermission{}, ErrNotFound
	} else if err != nil {
		return models.PathPermission{}, errors.Wrap(err, "Could not retrieve path permission")
	}

	return models.PathPermission{ID: id, PathPattern: pathPattern}, nil
}

func (db *db) RefreshPathPermissions(ctx context.Context, perms []models.PathPermission) error {
	c := db.pool.Get()
	defer c.Close()

	_ = c.Send("MULTI")
	_ = c.Send("DEL", pathPermissionHash)

	for _, p := range perms {
		_ = c.Send("HSET", pathPermissionHash, p.ID, p.PathPattern)
	}

	if _, err := c.Do("EXEC"); err != nil {
		return errors.Wrap(err, "Refresh path permissions")
	}
	return nil
}
