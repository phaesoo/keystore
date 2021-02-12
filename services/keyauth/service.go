package keyauth

import (
	"context"
	"encoding/json"

	"github.com/phaesoo/shield/pkg/db"
	"github.com/square/go-jose"
)

type Service struct {
	db *db.DB
}

func NewService(db *db.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Verify(ctx context.Context, tokenString, path, query string) error {

	// decode JWT token without verifying the signature

	token, err := jose.ParseSigned(tokenString)
	if err != nil {
		return err
	}

	b := token.UnsafePayloadWithoutVerification()

	payload := struct {
		AccessKey string `json:"access_key"`
		Nonce     string `json:"nonce"`
	}{}
	if err := json.Unmarshal(b, &payload); err != nil {
		return err
	}

	_, err = token.Verify([]byte("456"))
	if err != nil {
		return err
	}

	return nil
}
