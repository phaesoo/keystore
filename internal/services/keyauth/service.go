package keyauth

import (
	"context"
	"encoding/json"
	"log"

	"github.com/gobwas/glob"
	"github.com/phaesoo/shield/internal/models"
	"github.com/pkg/errors"
	"github.com/square/go-jose"
)

type serviceRepo interface {
	AuthKey(ctx context.Context, accessKey string) (models.AuthKey, error)
	PathPermission(ctx context.Context, id int) (models.PathPermission, error)
	PathPermissionIDs(ctx context.Context, keyID int) ([]int, error)
	RefreshPathPermissions(ctx context.Context) error
}

type Service struct {
	repo serviceRepo
}

func NewService(repo serviceRepo) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) Initialize(ctx context.Context) error {
	return s.repo.RefreshPathPermissions(ctx)
}

func (s *Service) Verify(ctx context.Context, token, urlPath, rawQuery string) (string, error) {
	signed, err := jose.ParseSigned(token)
	if err != nil {
		return "", err
	}

	// Decode JWT token without verifying the signature
	b := signed.UnsafePayloadWithoutVerification()
	var payload models.Payload
	if err := json.Unmarshal(b, &payload); err != nil {
		return "", err
	}

	authKey, err := s.repo.AuthKey(ctx, payload.AccessKey)
	if err != nil {
		return "", err
	}

	// TODO: Validate query with signature

	_, err = signed.Verify([]byte(authKey.SecretKey))
	if err != nil {
		log.Print("Verification failed")
		return "", err
	}

	permIDs, err := s.repo.PathPermissionIDs(ctx, authKey.ID)
	if err != nil {
		return "", err
	}

	for _, id := range permIDs {
		perm, err := s.repo.PathPermission(ctx, id)
		if err != nil {
			return "", err
		}
		g, err := glob.Compile(perm.PathPattern, '/')
		if err != nil {
			return "", err
		}
		if g.Match(urlPath) {
			return authKey.UserUUID, nil
		}
	}

	return "", errors.New("Url path not allowed")
}
