package mockrepo

import (
	"context"

	"github.com/phaesoo/shield/internal/models"
)

func (m *MockRepo) AuthKey(ctx context.Context, accessKey string) (models.AuthKey, error) {
	args := m.MethodCalled("AuthKey", ctx, accessKey)
	return args.Get(0).(models.AuthKey), args.Error(1)
}

func (m *MockRepo) PathPermission(ctx context.Context, id int) (models.PathPermission, error) {
	args := m.MethodCalled("PathPermission", ctx, id)
	return args.Get(0).(models.PathPermission), args.Error(1)
}

func (m *MockRepo) RefreshPathPermissions(ctx context.Context) error {
	args := m.MethodCalled("RefreshPathPermissions", ctx)
	return args.Error(0)
}

func (m *MockRepo) PathPermissionIDs(ctx context.Context, accessKey string) ([]int, error) {
	args := m.MethodCalled("PathPermissionIDs", ctx, accessKey)
	return args.Get(0).([]int), args.Error(1)
}
