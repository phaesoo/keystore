package db

import (
	"github.com/phaesoo/shield/configs"
	"github.com/stretchr/testify/suite"
)

const testDatabase string = "shield_test"

type TestSuite struct {
	suite.Suite
	Conn *DB
}

func (s *TestSuite) SetupSuite() {
	mc := configs.Get().Mysql
	mc.Database = testDatabase // Set test database
	db, err := NewDB("mysql", DSN(mc.User, mc.Password, mc.Database, mc.Host, mc.Port))
	if err != nil {
		panic(err)
	}
	s.Conn = db
}

func (s *TestSuite) Reset() {
	tx := s.Conn.MustBegin()
	tx.MustExec("SET FOREIGN_KEY_CHECKS=0;")
	tx.MustExec("TRUNCATE TABLE auth_key;")
	tx.MustExec("TRUNCATE TABLE path_permission;")
	tx.MustExec("TRUNCATE TABLE auth_key_path_permissions;")
	tx.MustExec("SET FOREIGN_KEY_CHECKS=1;")
	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

func (s *TestSuite) TearDownSuite() {
	s.Reset()
}
