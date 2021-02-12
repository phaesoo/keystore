package store

import (
	"context"
	"testing"
	"time"

	"github.com/phaesoo/shield/pkg/memdb"
	"github.com/stretchr/testify/assert"
)

const (
	testKeyA = "testKeyA"
	testKeyB = "testKeyB"
)

func Test_RedisCommon(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	pool := memdb.NewPool("0.0.0.0:6239", 1)
	defer pool.Close()
	t.Run("it sets and gets int correctly", func(t *testing.T) {
		r := NewStore(pool)
		expected := 1

		// set first
		err := r.Set(context.Background(), testKeyA, expected)
		assert := assert.New(t)
		assert.NoError(err)

		// get and check
		got, err := r.GetInt64(context.Background(), testKeyA)
		assert.NoError(err)
		assert.EqualValues(expected, got)
	})

	t.Run("it sets and gets string correctly", func(t *testing.T) {
		r := NewStore(pool)
		expected := "Hello World"
		// set first
		err := r.Set(context.Background(), testKeyB, expected)
		assert := assert.New(t)
		assert.NoError(err)

		// get and check
		got, err := r.GetString(context.Background(), testKeyB)
		assert.NoError(err)
		assert.EqualValues(expected, got)
	})

	t.Run("it runs IncrBy correctly", func(t *testing.T) {
		assert := assert.New(t)
		r := NewStore(pool)
		// set key value to 0
		err := r.Set(context.Background(), testKeyA, 0)
		assert.NoError(err)
		expected := 5

		// call IncrBy and check
		got, err := r.IncrBy(context.Background(), testKeyA, int64(expected))
		assert.NoError(err)
		assert.EqualValues(expected, got)
	})

	t.Run("it gets auto increment id correctly", func(t *testing.T) {
		assert := assert.New(t)
		r := NewStore(pool)
		// reset key to zero
		key := testKeyA
		err := r.Set(context.Background(), key, 0)
		assert.NoError(err)

		// check first id
		id, err := r.GetAutoIncrID(context.Background(), key)
		assert.NoError(err)
		assert.EqualValues(1, id)

		// check next id
		id, err = r.GetAutoIncrID(context.Background(), key)
		assert.NoError(err)
		assert.EqualValues(2, id)
	})

	t.Run("it runs SETNX correctly", func(t *testing.T) {
		r := NewStore(pool)

		key := "NON-EXIST-KEY"
		// should set
		expected := int64(1)
		err := r.SetNX(context.Background(), key, expected)
		assert := assert.New(t)
		assert.NoError(err)
		got, _ := r.GetInt64(context.Background(), key)
		assert.EqualValues(expected, got)

		// should not set
		expected = 2
		err = r.SetNX(context.Background(), key, expected)
		assert.NoError(err)
		got, _ = r.GetInt64(context.Background(), key)
		assert.NotEqual(expected, got)
	})

	t.Run("it runs EXISTS correctly", func(t *testing.T) {
		r := NewStore(pool)

		nonExistKey := "NON-EXIST-KEY2"
		existKey := "EXIST-KEY2"
		_ = r.Set(context.Background(), existKey, 10)

		tests := []struct {
			title    string
			key      string
			expected bool
		}{
			{"it returns false for the key not exists", nonExistKey, false},
			{"it returns true for the key exists", existKey, true},
		}

		for _, tt := range tests {
			t.Run(tt.title, func(t *testing.T) {
				got, err := r.Exists(context.Background(), tt.key)
				assert := assert.New(t)
				assert.Equal(tt.expected, got)
				assert.NoError(err)
			})
		}
	})

	t.Run("it sets value with expiration", func(t *testing.T) {
		assert := assert.New(t)
		r := NewStore(pool)

		key := "SetWithExpirationTestKey"
		expected := int64(10)

		err := r.SetWithExpiration(context.Background(), key, expected, 1)
		assert.NoError(err)
		got, err := r.GetInt64(context.Background(), key)
		assert.Equal(expected, got)
		assert.NoError(err)
		time.Sleep(2 * time.Second)
		_, err = r.GetInt64(context.Background(), key)
		assert.Error(err)

	})
}
