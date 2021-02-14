package db

import (
	"context"
	"fmt"

	"github.com/phaesoo/shield/internal/models"
)

func (db *DB) AuthKey(ctx context.Context, accessKey string) (models.AuthKey, error) {
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
func (db *DB) PathPermissions(ctx context.Context) ([]models.PathPermission, error) {
	perms := []models.PathPermission{}

	rows, err := db.conn.Queryx(`SELECT id, path_pattern FROM path_permission`)
	if err != nil {
		return perms, err
	}

	for rows.Next() {
		perm := struct {
			ID          int    `db:"id"`
			PathPattern string `db:"path_pattern"`
		}{}

		err = rows.StructScan(&perm)
		if err != nil {
			return perms, err
		}
		perms = append(perms, models.PathPermission{
			ID:          perm.ID,
			PathPattern: perm.PathPattern,
		})
	}
	return perms, nil
}

func (db *DB) PathPermissionIDs(ctx context.Context, keyID int) ([]int, error) {
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
