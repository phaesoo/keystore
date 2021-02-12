package server

import (
	"context"
	"expvar"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/phaesoo/shield/internal/resp"
)

type Server struct {
	server *http.Server
	router *chi.Mux
	config ServerConfig
}

type ServerConfig struct {
	Profile bool
	Metrics bool
}

func NewServer(address string, config ServerConfig) *Server {
	router := chi.NewRouter()

	httpServer := &http.Server{
		Addr:         address,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	s := &Server{
		server: httpServer,
		router: router,
		config: config,
	}

	s.configureMiddleware()
	s.configureRoutes()

	return s
}

func (s *Server) handleIndex(w http.ResponseWriter, r *http.Request) {
	log.Print("Call health check")
	resp.OK(w, nil, "Yes, Boss.")
}

func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	resp.Fail(w, http.StatusNotFound, nil, "The requested endpoint could not be found", resp.CodeUnknownEndpoint)
}

// configureRoutes contains all the management routes that the server handles. Keep in mind that
// services register their own routes.
func (s *Server) configureRoutes() {
	// Enable profiling routes. See: https://golang.org/pkg/net/http/pprof/
	if s.config.Profile {
		s.router.Mount("/debug", middleware.Profiler())
	}

	// Enable expvar metrics. See: https://golang.org/pkg/expvar/
	if s.config.Metrics {
		s.router.Method(http.MethodGet, "/debug/vars", expvar.Handler())
	}

	// Health Check
	s.router.Get("/", s.handleIndex)

	s.router.NotFound(s.handleNotFound)
}

func (s *Server) configureMiddleware() {
}

// Router returns server router
func (s *Server) Router() *chi.Mux {
	return s.router
}

// Listen starts the server on the address given in the server configuration.
func (s *Server) Listen() error {
	s.printRoutes()
	log.Printf("Start to listen: %s", s.server.Addr)
	return s.server.ListenAndServe()
}

// Shutdown the server gracefully with a 30 second timeout
func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) printRoutes() {
	routes := []string{}
	walkRoute := func(method, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.ReplaceAll(route, "/*/", "/")
		routes = append(routes, fmt.Sprintf("%s %s", method, route))
		return nil
	}
	_ = chi.Walk(s.router, walkRoute)
	log.Printf("Serving routes: \n%s", strings.Join(routes, "\n"))
}
