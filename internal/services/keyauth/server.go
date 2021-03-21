package keyauth

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/phaesoo/shield/internal/resp"
)

const (
	routeStatus = "/token/verify"
)

type Server struct {
	service *Service
}

func NewServer(service *Service) *Server {
	return &Server{
		service: service,
	}
}

func (s *Server) RegisterRoutes(router *chi.Mux) {
	router.Put(routeStatus, s.HandleTokenVerify)
}

func (s *Server) HandleTokenVerify(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("token")
	if token != "" {
		resp.BadRequest(w, resp.CodeInvalidRequest)
		return
	}

	log.Print("HandlePostVerify called")

	if _, err := s.service.Verify(context.Background(), token, "", ""); err != nil {
		resp.Error(w, err, "")
		return
	}
	resp.OK(w, nil, "OK")
}
