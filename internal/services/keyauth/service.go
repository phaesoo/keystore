package keyauth

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gobwas/glob"
	"github.com/phaesoo/shield/internal/models"
	"github.com/phaesoo/shield/pkg/db"
	"github.com/pkg/errors"
	"github.com/square/go-jose"
)

type serviceStore interface {
	PathPermission(ctx context.Context, id int) (models.PathPermission, error)
	RefreshPathPermissions(ctx context.Context, perms []models.PathPermission) error
}

type serviceRepo interface {
	AuthKey(ctx context.Context, accessKey string) (models.AuthKey, error)
	PathPermissionIDs(ctx context.Context, keyID int) ([]int, error)
}

type Service struct {
	store serviceStore
	repo  serviceRepo
	db    *db.DB
}

func NewService(store serviceStore, db *db.DB) *Service {
	return &Service{
		store: store,
		db:    db,
	}
}

func (s *Service) Initialize(ctx context.Context) error {
	perms := []models.PathPermission{}

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
		perms = append(perms, models.PathPermission{
			ID:          perm.ID,
			PathPattern: perm.PathPattern,
		})
	}

	if err := s.store.RefreshPathPermissions(context.Background(), perms); err != nil {
		return err
	}

	return nil
}

func (s *Service) Verify(ctx context.Context, token, urlPath, rawQuery string) error {
	signed, err := jose.ParseSigned(token)
	if err != nil {
		return err
	}

	// Decode JWT token without verifying the signature
	b := signed.UnsafePayloadWithoutVerification()
	payload := struct {
		AccessKey string `json:"access_key"`
		Nonce     string `json:"nonce"`
		Signature string `json:"signature"`
	}{}
	if err := json.Unmarshal(b, &payload); err != nil {
		return err
	}

	// TODO: Validate query with signature

	authKey, err := s.repo.AuthKey(ctx, payload.AccessKey)
	if err != nil {
		return err
	}

	_, err = signed.Verify([]byte(authKey.SecretKey))
	if err != nil {
		log.Print("Verification failed")
		return err
	}

	permIDs, err := s.repo.PathPermissionIDs(ctx, authKey.ID)
	if err != nil {
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
