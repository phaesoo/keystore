package db

import (
	"context"
	"testing"

	"github.com/bxcodec/faker/v3"
	"github.com/phaesoo/shield/internal/models"
	"github.com/phaesoo/shield/pkg/db"
	"github.com/stretchr/testify/suite"
)

func TestKeyAuth(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration tests")
	}
	suite.Run(t, new(KeyAuthTestSuite))
}

type KeyAuthTestSuite struct {
	db.TestSuite
	db *DB
}

func (ts *KeyAuthTestSuite) SetupSuite() {
	ts.TestSuite.SetupSuite()
	ts.db = NewDB(ts.Conn)
}

func (ts *KeyAuthTestSuite) Test_SetAuthKey() {
	authKey := models.AuthKey{}
	ts.NoError(faker.FakeData(&authKey))

	err := ts.db.SetAuthKey(context.Background(), authKey)
	ts.NoError(err)
}
