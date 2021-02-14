package keyauth

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gobwas/glob"
	"github.com/phaesoo/shield/internal/models"
	"github.com/phaesoo/shield/internal/store"
	"github.com/phaesoo/shield/pkg/db"
	"github.com/pkg/errors"
	"github.com/square/go-jose"
)

type serviceStore interface {
	AuthKey(ctx context.Context, accessKey string) (models.AuthKey, error)
	SetAuthKey(ctx context.Context, authKey models.AuthKey) error
	PathPermission(ctx context.Context, id int) (models.PathPermission, error)
	RefreshPathPermissions(ctx context.Context, perms []models.PathPermission) error
}

type serviceRepo interface {
	AuthKey(ctx context.Context, accessKey string) (models.AuthKey, error)
	PathPermissionIDs(ctx context.Context, keyID int) ([]int, error)
	PathPermissions(ctx context.Context) ([]models.PathPermission, error)
}

type Service struct {
	repo  serviceRepo
	store serviceStore
}

func NewService(repo serviceRepo, store serviceStore, db *db.DB) *Service {
	return &Service{
		repo:  repo,
		store: store,
	}
}

func (s *Service) Initialize(ctx context.Context) error {
	perms, err := s.repo.PathPermissions(ctx)
	if err != nil {
		return err
	}

	if err := s.store.RefreshPathPermissions(context.Background(), perms); err != nil {
		return err
	}

	return nil
}

func (s *Service) findAuthKey(ctx context.Context, accessKey string) (models.AuthKey, error) {
	var authKey models.AuthKey
	var err error
	authKey, err = s.store.AuthKey(ctx, accessKey)
	if err != nil {
		if err != store.ErrNotFound {
			return authKey, err
		}
		authKey, err = s.repo.AuthKey(ctx, accessKey)
		if err != nil {
			return authKey, err
		}
		if err := s.store.SetAuthKey(ctx, authKey); err != nil {
			return authKey, err
		}
	}
	return authKey, nil
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

	authKey, err := s.findAuthKey(ctx, payload.AccessKey)
	if err != nil {
		return err
	}

	// TODO: Validate query with signature

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
