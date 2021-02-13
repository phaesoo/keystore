package store

import (
	"context"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

const (
	pathPermissionKey = "path-permissions"
)

type pathPermStore interface {
	PathPermission(ctx context.Context, id int) (PathPermission, error)
	RefreshPathPermissions(ctx context.Context, perms []PathPermission) error
}

func (db *db) PathPermission(ctx context.Context, id int) (PathPermission, error) {
	c := db.pool.Get()
	defer c.Close()

	pathPattern, err := redis.String(c.Do("HGET", pathPermissionKey, id))
	if err == redis.ErrNil {
		return PathPermission{}, ErrNotFound
	} else if err != nil {
		return PathPermission{}, errors.Wrap(err, "Could not retrieve path permission")
	}

	return PathPermission{ID: id, PathPattern: pathPattern}, nil
}

func (db *db) RefreshPathPermissions(ctx context.Context, perms []PathPermission) error {
	c := db.pool.Get()
	defer c.Close()

	_ = c.Send("MULTI")
	_ = c.Send("DEL", pathPermissionKey)

	for _, p := range perms {
		_ = c.Send("HSET", pathPermissionKey, p.ID, p.PathPattern)
	}

	if _, err := c.Do("EXEC"); err != nil {
		return errors.Wrap(err, "Refresh path permissions")
	}
	return nil
}
