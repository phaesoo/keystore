package repo

import (
	"context"

	"github.com/phaesoo/shield/internal/models"
	"github.com/phaesoo/shield/internal/repo/cache"
)

type keyauthRepo interface {
	AuthKey(ctx context.Context, accessKey string) (models.AuthKey, error)
	PathPermission(ctx context.Context, id int) (models.PathPermission, error)
	RefreshPathPermissions(ctx context.Context) error
	PathPermissionIDs(ctx context.Context, keyID int) ([]int, error)
}

func (r *repo) AuthKey(ctx context.Context, accessKey string) (models.AuthKey, error) {
	var authKey models.AuthKey
	var err error
	authKey, err = r.cache.AuthKey(ctx, accessKey)
	if err != nil {
		if err != cache.ErrNotFound {
			return authKey, err
		}
		authKey, err = r.db.AuthKey(ctx, accessKey)
		if err != nil {
			return authKey, err
		}
		if err := r.cache.SetAuthKey(ctx, authKey); err != nil {
			return authKey, err
		}
	}
	return authKey, nil
}

func (r *repo) PathPermission(ctx context.Context, id int) (models.PathPermission, error) {
	return r.cache.PathPermission(ctx, id)
}

func (r *repo) RefreshPathPermissions(ctx context.Context) error {
	perms, err := r.db.PathPermissions(ctx)
	if err != nil {
		return err
	}
	return r.cache.RefreshPathPermissions(ctx, perms)
}

func (r *repo) PathPermissionIDs(ctx context.Context, keyID int) ([]int, error) {
	var permIDs []int
	var err error
	permIDs, err = r.cache.PathPermissionIDs(ctx, keyID)
	if err != nil {
		if err != cache.ErrNotFound {
			return permIDs, err
		}
		permIDs, err = r.db.PathPermissionIDs(ctx, keyID)
		if err != nil {
			return permIDs, err
		}
		if err := r.cache.SetPathPermissionIDs(ctx, keyID, permIDs); err != nil {
			return permIDs, err
		}
	}
	return permIDs, nil
}
