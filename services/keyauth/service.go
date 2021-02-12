package keyauth

import (
	"context"

	"github.com/phaesoo/shield/pkg/db"
)

type Service struct {
	db *db.DB
}

func NewService(db *db.DB) *Service {
	return &Service{db: db}
}

func (s *Service) Verify(ctx context.Context) error {
	return nil
}
