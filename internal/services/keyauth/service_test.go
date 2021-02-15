package keyauth

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/phaesoo/shield/internal/models"
	"github.com/phaesoo/shield/internal/repo/mockrepo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestService_Verify(t *testing.T) {
	expectedUserUUID := "uuid-1"
	expectedPattern := "/markets/all"

	t.Run("It runs without error and returns expected user uuid", func(t *testing.T) {
		accessKey := uuid.NewString()
		secretKey := uuid.NewString()

		payload := models.Payload{
			AccessKey: accessKey,
			Nonce:     uuid.NewString(),
			Signature: "1",
		}
		token, err := payload.Encrypt(secretKey)
		assert.NoError(t, err)

		repo := mockrepo.NewMockRepo()

		expectedAuthKey := models.AuthKey{
			ID:        1,
			AccessKey: accessKey,
			SecretKey: secretKey,
			UserUUID:  expectedUserUUID,
		}
		repo.On("AuthKey", mock.Anything, mock.Anything).Return(
			expectedAuthKey,
			nil,
		).Once()

		expectedPathPermissionIDs := []int{1}
		repo.On("PathPermissionIDs", mock.Anything, expectedAuthKey.AccessKey).Return(
			expectedPathPermissionIDs,
			nil,
		).Once()

		repo.On("PathPermission", mock.Anything, expectedPathPermissionIDs[0]).Return(
			models.PathPermission{
				ID:          expectedPathPermissionIDs[0],
				PathPattern: expectedPattern,
			},
			nil,
		).Once()

		service := NewService(repo)

		userUUID, err := service.Verify(context.Background(), token, expectedPattern, "1")
		assert.Equal(t, expectedUserUUID, userUUID)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})
	t.Run("It returns error if secret key is not equal", func(t *testing.T) {
		accessKey := uuid.NewString()
		secretKey := uuid.NewString()
		secretKey2 := uuid.NewString()

		payload := models.Payload{
			AccessKey: accessKey,
			Nonce:     uuid.NewString(),
			Signature: "1",
		}
		token, err := payload.Encrypt(secretKey)
		assert.NoError(t, err)

		repo := mockrepo.NewMockRepo()

		expectedAuthKey := models.AuthKey{
			ID:        1,
			AccessKey: accessKey,
			SecretKey: secretKey2,
			UserUUID:  expectedUserUUID,
		}
		repo.On("AuthKey", mock.Anything, mock.Anything).Return(
			expectedAuthKey,
			nil,
		).Once()

		service := NewService(repo)

		userUUID, err := service.Verify(context.Background(), token, expectedPattern, "1")
		assert.Equal(t, "", userUUID)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
	t.Run("It returns error if permission ids are empty", func(t *testing.T) {
		accessKey := uuid.NewString()
		secretKey := uuid.NewString()

		payload := models.Payload{
			AccessKey: accessKey,
			Nonce:     uuid.NewString(),
			Signature: "1",
		}
		token, err := payload.Encrypt(secretKey)
		assert.NoError(t, err)

		repo := mockrepo.NewMockRepo()

		expectedAuthKey := models.AuthKey{
			ID:        1,
			AccessKey: accessKey,
			SecretKey: secretKey,
			UserUUID:  expectedUserUUID,
		}
		repo.On("AuthKey", mock.Anything, mock.Anything).Return(
			expectedAuthKey,
			nil,
		).Once()

		repo.On("PathPermissionIDs", mock.Anything, expectedAuthKey.AccessKey).Return(
			[]int{},
			nil,
		).Once()

		service := NewService(repo)

		userUUID, err := service.Verify(context.Background(), token, expectedPattern, "1")
		assert.Equal(t, "", userUUID)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
	t.Run("It returns error if there is no matched pattern", func(t *testing.T) {
		accessKey := uuid.NewString()
		secretKey := uuid.NewString()

		payload := models.Payload{
			AccessKey: accessKey,
			Nonce:     uuid.NewString(),
			Signature: "1",
		}
		token, err := payload.Encrypt(secretKey)
		assert.NoError(t, err)

		repo := mockrepo.NewMockRepo()

		expectedAuthKey := models.AuthKey{
			ID:        1,
			AccessKey: accessKey,
			SecretKey: secretKey,
			UserUUID:  expectedUserUUID,
		}
		repo.On("AuthKey", mock.Anything, mock.Anything).Return(
			expectedAuthKey,
			nil,
		).Once()

		repo.On("PathPermissionIDs", mock.Anything, expectedAuthKey.AccessKey).Return(
			[]int{},
			nil,
		).Once()

		service := NewService(repo)

		userUUID, err := service.Verify(context.Background(), token, "/markets/all", "1")
		assert.Equal(t, "", userUUID)
		assert.Error(t, err)
		repo.AssertExpectations(t)
	})
}
