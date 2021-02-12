package request

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBindJSON(t *testing.T) {
	t.Run("Does not allow non-JSON content types", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/", nil)
		r.Header.Set("Content-Type", "text/plain")

		err := BindJSON(r, nil)
		assert.Equal(t, ErrInvalidContentType, err)
	})

	t.Run("Does not allow trailing body data", func(t *testing.T) {
		body := []byte("{\"id\": \"123\"}=true")
		buff := bytes.NewBuffer(body)
		r := httptest.NewRequest(http.MethodPost, "/", buff)
		r.Header.Set("Content-Type", "application/json")

		var b struct {
			ID string `json:"id"`
		}
		err := BindJSON(r, &b)
		assert.Equal(t, ErrExpectedEOF, err)
	})

	t.Run("Can parse a valid JSON request", func(t *testing.T) {
		body := []byte("{\"id\": \"test_id\"}")
		buff := bytes.NewBuffer(body)
		r := httptest.NewRequest(http.MethodPost, "/", buff)
		r.Header.Set("Content-Type", "application/json; charset=utf-8")

		var b struct {
			ID string `json:"id"`
		}
		err := BindJSON(r, &b)

		assert := assert.New(t)
		assert.NoError(err)
		assert.Equal("test_id", b.ID)
	})
}
