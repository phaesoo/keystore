package cache

import (
	"testing"

	"github.com/bxcodec/faker/v3"

	"github.com/phaesoo/shield/internal/models"
	"github.com/phaesoo/shield/pkg/memdb"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestKeyAuth(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration tests")
	}
	suite.Run(t, new(KeyAuthTestSuite))
}

func Test_authKeyHash(t *testing.T) {
	assert := assert.New(t)

	var accessKey string
	assert.NoError(faker.FakeData(&accessKey))

	res := authKeyHash(accessKey)
	assert.Equal(res, authKeyHashPrefix+accessKey)
}

type KeyAuthTestSuite struct {
	memdb.TestSuite
	cache *Cache
}

func (ts *KeyAuthTestSuite) SetupSuite() {
	ts.TestSuite.SetupSuite()
	ts.cache = NewCache(ts.Pool)
}

func (ts *KeyAuthTestSuite) Test_SetAuthKey() {
	var authKey models.AuthKey
	ts.NoError(faker.FakeData(&authKey))

	err := ts.cache.SetAuthKey(authKey)
	ts.NoError(err)
}

func (ts *KeyAuthTestSuite) Test_GetAuthKey() {
	var authKey models.AuthKey
	ts.NoError(faker.FakeData(&authKey))

	err := ts.cache.SetAuthKey(authKey)
	ts.NoError(err)

	res, err := ts.cache.AuthKey(authKey.AccessKey)
	ts.NoError(err)
	ts.EqualValues(authKey, res)
}

func (ts *KeyAuthTestSuite) Test_RefreshPathPermissions() {
	perms := make([]models.PathPermission, 10)
	for i, _ := range perms {
		perms[i].ID = i
		faker.FakeData(&perms[i].PathPattern)
	}

	err := ts.cache.RefreshPathPermissions(perms)
	ts.NoError(err)
}

func (ts *KeyAuthTestSuite) Test_PathPermission() {
	perms := make([]models.PathPermission, 20)
	for i, _ := range perms {
		perms[i].ID = i
		faker.FakeData(&perms[i].PathPattern)
	}

	err := ts.cache.RefreshPathPermissions(perms)
	ts.NoError(err)

	for _, p := range perms {
		res, err := ts.cache.PathPermission(p.ID)
		ts.NoError(err)
		ts.EqualValues(p, res)
	}
}

func (ts *KeyAuthTestSuite) Test_SetPathPermissionIDs() {
	var accessKey string
	ts.NoError(faker.FakeData(&accessKey))

	var permIDs []int
	ts.NoError(faker.FakeData(&permIDs))

	err := ts.cache.SetPathPermissionIDs(accessKey, permIDs)
	ts.NoError(err)
}

func (ts *KeyAuthTestSuite) Test_PathPermissionIDs() {
	var accessKey string
	ts.NoError(faker.FakeData(&accessKey))

	var permIDs []int
	ts.NoError(faker.FakeData(&permIDs))

	err := ts.cache.SetPathPermissionIDs(accessKey, permIDs)
	ts.NoError(err)

	res, err := ts.cache.PathPermissionIDs(accessKey)
	ts.NoError(err)
	ts.EqualValues(permIDs, res)
}
