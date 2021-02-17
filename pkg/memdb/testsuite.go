package memdb

import (
	"github.com/gomodule/redigo/redis"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
	Pool *redis.Pool
	Conn redis.Conn
}

func (s *TestSuite) ClearKeys(keys ...string) {
	for _, key := range keys {
		_, err := s.Conn.Do("DEL", key)
		s.NoError(err)
	}
}

func (s *TestSuite) SetupSuite() {
	s.Pool = NewTestPool()
	s.Conn = s.Pool.Get()
}

func (s *TestSuite) TearDownSuite() {
	s.NoError(s.Conn.Flush())
	s.NoError(s.Conn.Close())
	s.NoError(s.Pool.Close())
}
