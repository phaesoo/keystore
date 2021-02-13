package repo

import (
	"context"
	"fmt"

	"github.com/phaesoo/shield/internal/models"
)

type keyauthRepo interface {
	PathPermissionIDs(ctx context.Context, id int) (models.PathPermission, error)
}

func (db *db) AuthKey(ctx context.Context, accessKey string) (models.AuthKey, error) {
	k := struct {
		ID        int    `db:"id"`
		AccessKey string `db:"access_key"`
		SecretKey string `db:"secret_key"`
		UserUUID  string `db:"user_uuid"`
	}{}

	if err := db.conn.Get(&k, fmt.Sprintf(`
		SELECT *
		FROM auth_key
		WHERE access_key = %s
		`, accessKey)); err != nil {
		return models.AuthKey{}, err
	}

	return models.AuthKey{
		ID:        k.ID,
		AccessKey: k.AccessKey,
		SecretKey: k.SecretKey,
		UserUUID:  k.UserUUID,
	}, nil
}

func (db *db) PathPermissionIDs(ctx context.Context, keyID int) ([]int, error) {
	var output []int
	if err := db.conn.Select(&output, fmt.Sprintf(`
		SELECT B.permission_id
		FROM auth_key A
		JOIN auth_key_path_permissions B on A.id = B.key_id
		WHERE A.id = %d
		`, keyID)); err != nil {
		return output, err
	}
	return output, nil
}
