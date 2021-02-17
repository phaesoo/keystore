package db

import (
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
	var authKey models.AuthKey
	ts.NoError(faker.FakeData(&authKey))

	err := ts.db.SetAuthKey(authKey)
	ts.NoError(err)
}

func (ts *KeyAuthTestSuite) Test_AuthKey() {
	var authKey models.AuthKey
	ts.NoError(faker.FakeData(&authKey))

	err := ts.db.SetAuthKey(authKey)
	ts.NoError(err)

	ts.Run("It returns expected object", func() {
		res, err := ts.db.AuthKey(authKey.AccessKey)
		ts.NoError(err)
		ts.EqualValues(authKey, res)
	})
	ts.Run("It returns error with unknown access key", func() {
		var accessKey string
		ts.NoError(faker.FakeData(&accessKey))

		res, err := ts.db.AuthKey(accessKey)
		ts.Error(err)
		ts.EqualValues(models.AuthKey{}, res)
	})
}

func (ts *KeyAuthTestSuite) Test_PathPermissions() {
	var authKey models.AuthKey
	ts.NoError(faker.FakeData(&authKey))

	err := ts.db.SetAuthKey(authKey)
	ts.NoError(err)

	ts.Run("It returns expected object", func() {
		res, err := ts.db.AuthKey(authKey.AccessKey)
		ts.NoError(err)
		ts.EqualValues(authKey, res)
	})
	ts.Run("It returns error with unknown access key", func() {
		var accessKey string
		ts.NoError(faker.FakeData(&accessKey))

		res, err := ts.db.AuthKey(accessKey)
		ts.Error(err)
		ts.EqualValues(models.AuthKey{}, res)
	})
}
