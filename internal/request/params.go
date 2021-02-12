package request

import (
	"net/http"

	"github.com/go-chi/chi"
)


// GetParam is wrapper function of chi URLParam
func GetParam(r* http.Request, key string) string {
	return chi.URLParam(r, key)
}