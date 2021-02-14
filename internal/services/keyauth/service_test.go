package keyauth

import (
	"context"
	"testing"

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

	err := s.service.Verify(context.Background(), tokenString, "temp", "1")
	s.NoError(err)
	//err := s.service.GenerateSchedule(context.Background(), "20201010")
	//s.NoError(err)
}
