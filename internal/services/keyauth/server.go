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
	router.Put(routeStatus, s.HandlePostVerify)
}

func (s *Server) HandlePostVerify(w http.ResponseWriter, r *http.Request) {
	log.Print("HandlePostVerify called")

	if _, err := s.service.Verify(context.Background(), "", "", ""); err != nil {
		resp.Error(w, err, "")
		return
	}
	resp.OK(w, nil, "OK")
}
