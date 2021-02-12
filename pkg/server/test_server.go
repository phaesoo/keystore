package server

import (
	"net/http"
	"net/http/httptest"

	"github.com/go-chi/chi"
)

type TestServer struct {
	http.Handler
}

type RoutedService interface {
	RegisterRoutes(*chi.Mux)
}

// NewTestServer creates a new test server based on the provided routed service
func NewTestServer(service RoutedService) *TestServer {
	router := chi.NewRouter()
	service.RegisterRoutes(router)
	return &TestServer{router}
}

func (t *TestServer) Request(r *http.Request) *httptest.ResponseRecorder {
	w := httptest.NewRecorder()
	t.ServeHTTP(w, r)
	return w
}
