package request

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_validateContentType(t *testing.T) {
	t.Run("Test no error", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte{}))
		r.Header.Set("Content-Type", "application/json")

		err := validateContentType(r, "application/json")
		assert.NoError(t, err)
	})

	t.Run("Test ErrInvalidContentType", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer([]byte{}))
		r.Header.Set("Content-Type", "text/plain")

		err := validateContentType(r, "application/json")
		assert.Equal(t, ErrInvalidContentType, err)
	})
}
