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
	PathPermissionIDs(ctx context.Context, accessKey string) ([]int, error)
	RefreshPathPermissions(ctx context.Context) error
}

// Service is implementation of keyauth service
type Service struct {
	repo serviceRepo
}

// NewService creates keyauth service
func NewService(repo serviceRepo) *Service {
	return &Service{
		repo: repo,
	}
}

// Initialize preset for service
func (s *Service) Initialize(ctx context.Context) error {
	return s.repo.RefreshPathPermissions(ctx)
}

// Verify JWT(JWS) token and returns user uuid
func (s *Service) Verify(ctx context.Context, token, urlPath, queryString string) (string, error) {
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

	if err := payload.Validate(queryString); err != nil {
		return "", err
	}

	authKey, err := s.repo.AuthKey(ctx, payload.AccessKey)
	if err != nil {
		return "", err
	}

	_, err = signed.Verify([]byte(authKey.SecretKey))
	if err != nil {
		log.Print("Verification failed")
		return "", err
	}

	permIDs, err := s.repo.PathPermissionIDs(ctx, authKey.AccessKey)
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
