package keyauth

import (
	"context"
	"testing"

	"github.com/phaesoo/shield/internal/models"
	"github.com/phaesoo/shield/internal/repo/mockrepo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestService(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration test")
	}
	suite.Run(t, new(ServiceTestSuite))
}

type ServiceTestSuite struct {
	suite.Suite
	service *Service
}

func (s *ServiceTestSuite) SetupSuite() {
	// var mc configs.MysqlConfig = configs.Get().Mysql
	// db, err := db.NewDB("mysql", db.DSN(mc.User, mc.Password, mc.Database, mc.Host, mc.Port))
	// if err != nil {
	// 	s.NoError(err)
	// }
	// s.service = &Service{db: db}
}

func (s *ServiceTestSuite) TearDownSuite() {
}

func (s *ServiceTestSuite) TestVerify() {
	tokenString := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhY2Nlc3Nfa2V5IjoiMTIzIiwibm9uY2UiOiI1ODA5YWVhYS0yZDhjLTQyZmEtOTk5Yi1iOTdmNjBhNTQ5YjQifQ.B3IYg6VvANcPjdKJZRlOrR2tFH2snpIA0pTYEiyFVuI"

	_, err := s.service.Verify(context.Background(), tokenString, "temp", "1")
	s.NoError(err)
	//err := s.service.GenerateSchedule(context.Background(), "20201010")
	//s.NoError(err)
}

func TestService_Verify(t *testing.T) {
	tokenString := "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJhY2Nlc3Nfa2V5IjoiMTIzIiwibm9uY2UiOiI1ODA5YWVhYS0yZDhjLTQyZmEtOTk5Yi1iOTdmNjBhNTQ5YjQifQ.B3IYg6VvANcPjdKJZRlOrR2tFH2snpIA0pTYEiyFVuI"

	repo := mockrepo.NewMockRepo()
	repo.On("AuthKey", mock.Anything, mock.Anything).Return(
		models.AuthKey{
			ID:        1,
			AccessKey: "123",
			SecretKey: "456",
			UserUUID:  "uuid-1",
		},
		nil,
	)
	repo.On("PathPermissionIDs", mock.Anything, mock.Anything).Return(
		[]int{1},
		nil,
	)

	repo.On("PathPermission", mock.Anything, mock.Anything).Return(
		models.PathPermission{
			ID:          1,
			PathPattern: "/markets/all",
		},
		nil,
	)

	service := NewService(repo)

	t.Run("", func(t *testing.T) {
		userUUID, err := service.Verify(context.Background(), tokenString, "/markets/all", "1")
		assert.Equal(t, "uuid-1", userUUID)
		assert.NoError(t, err)
	})
}
