package keyauth

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/gobwas/glob"
	"github.com/phaesoo/shield/internal/store"
	"github.com/phaesoo/shield/pkg/db"
	"github.com/pkg/errors"
	"github.com/square/go-jose"
)

type serviceStore interface {
	PathPermission(ctx context.Context, id int) (store.PathPermission, error)
	RefreshPathPermissions(ctx context.Context, perms []store.PathPermission) error
}

type Service struct {
	store serviceStore
	db    *db.DB
}

func NewService(store serviceStore, db *db.DB) *Service {
	return &Service{
		store: store,
		db:    db,
	}
}

func (s *Service) Initialize(ctx context.Context) error {
	perms := []store.PathPermission{}

	rows, err := s.db.Queryx(`SELECT id, path_pattern FROM path_permission`)
	if err != nil {
		return err
	}

	for rows.Next() {
		perm := struct {
			ID          int    `db:"id"`
			PathPattern string `db:"path_pattern"`
		}{}

		err = rows.StructScan(&perm)
		if err != nil {
			return err
		}
		perms = append(perms, store.PathPermission{
			ID:          perm.ID,
			PathPattern: perm.PathPattern,
		})
	}

	if err := s.store.RefreshPathPermissions(context.Background(), perms); err != nil {
		return err
	}

	return nil
}

func (s *Service) Verify(ctx context.Context, tokenString, urlPath, rawQuery string) error {
	token, err := jose.ParseSigned(tokenString)
	if err != nil {
		return err
	}

	// Decode JWT token without verifying the signature
	b := token.UnsafePayloadWithoutVerification()
	payload := struct {
		AccessKey string `json:"access_key"`
		Nonce     string `json:"nonce"`
		Signature string `json:"signature"`
	}{}
	if err := json.Unmarshal(b, &payload); err != nil {
		return err
	}

	// TODO: Validate query with signature

	authKey := struct {
		ID        int    `db:"id"`
		SecretKey string `db:"secret_key"`
		UserUUID  string `db:"user_uuid"`
	}{}

	if err := s.db.Get(&authKey, fmt.Sprintf(`
		SELECT id, secret_key, user_uuid
		FROM auth_key
		WHERE access_key = %s
		`, payload.AccessKey)); err != nil {
		return err
	}

	_, err = token.Verify([]byte(authKey.SecretKey))
	if err != nil {
		log.Print("Verification failed")
		return err
	}

	permIDs := []int{}
	if err := s.db.Select(&permIDs, fmt.Sprintf(`
		SELECT B.permission_id
		FROM auth_key A
		JOIN auth_key_path_permissions B on A.id = B.key_id
		WHERE A.id = %d
		`, authKey.ID)); err != nil {
		return err
	}

	for _, id := range permIDs {
		perm, err := s.store.PathPermission(context.Background(), id)
		if err != nil {
			return err
		}

		g, err := glob.Compile(perm.PathPattern, '/')
		if err != nil {
			return err
		}

		if g.Match(urlPath) {
			return nil
		}
	}

	return errors.Wrap(err, "url not allowed")
}
