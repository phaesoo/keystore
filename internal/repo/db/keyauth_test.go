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

	ts.Run("It returns expected objects", func() {
		defer ts.Reset()

		perms := make([]models.PathPermission, 2)
		tx := ts.Conn.MustBegin()
		for i := 0; i < 2; i++ {
			var pathPattern string
			ts.NoError(faker.FakeData(&pathPattern))
			pp := models.PathPermission{ID: i + 1, PathPattern: pathPattern}
			tx.MustExec(`
			INSERT INTO path_permission (path_pattern)
			VALUES (?)
			`, pp.PathPattern)

			perms[i] = pp
		}
		ts.NoError(tx.Commit())

		var accessKey string
		ts.NoError(faker.FakeData(&accessKey))

		res, err := ts.db.PathPermissions()
		ts.NoError(err)
		ts.EqualValues(perms, res)
	})
	ts.Run("It returns error if result is empty", func() {
		defer ts.Reset()

		var accessKey string
		ts.NoError(faker.FakeData(&accessKey))

		res, err := ts.db.PathPermissions()
		ts.Error(err)
		ts.Empty(res)
	})
}
