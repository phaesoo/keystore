package cache

import (
	"context"
	"testing"

	"github.com/bxcodec/faker/v3"

	"github.com/phaesoo/shield/internal/models"
	"github.com/phaesoo/shield/pkg/memdb"
	"github.com/stretchr/testify/suite"
)

func TestKeyAuth(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip integration tests")
	}
	suite.Run(t, new(KeyAuthTestSuite))
}

type KeyAuthTestSuite struct {
	memdb.TestSuite
	cache *Cache
}

func (ts *KeyAuthTestSuite) SetupSuite() {
	ts.TestSuite.SetupSuite()
	ts.cache = NewCache(ts.Pool)
}

func (ts *KeyAuthTestSuite) TestSetAuthKey() {
	authKey := models.AuthKey{}
	ts.NoError(faker.FakeData(&authKey))

	err := ts.cache.SetAuthKey(context.Background(), authKey)
	ts.NoError(err)
}

func (ts *KeyAuthTestSuite) TestGetAuthKey() {
	authKey := models.AuthKey{}
	ts.NoError(faker.FakeData(&authKey))

	err := ts.cache.SetAuthKey(context.Background(), authKey)
	ts.NoError(err)

	res, err := ts.cache.AuthKey(context.Background(), authKey.AccessKey)
	ts.NoError(err)

	ts.EqualValues(authKey, res)
}
