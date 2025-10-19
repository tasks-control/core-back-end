package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	v1 "github.com/tasks-control/core-back-end/api/v1"
	"github.com/tasks-control/core-back-end/internal/middleware"
	"github.com/tasks-control/core-back-end/internal/service"
	"github.com/tasks-control/core-back-end/pkg/utils"
)

type HTTPHandlers interface {
	HealthHandlers
	v1.ServerInterface
}

type HealthHandlers interface {
	Readiness(w http.ResponseWriter, r *http.Request)
}

type Server struct {
	httpServer *http.Server
	handlers   HTTPHandlers
	service    *service.Service
}

func NewServer(handlers HTTPHandlers, svc *service.Service) *Server {
	return &Server{
		handlers: handlers,
		service:  svc,
	}
}

func (s *Server) Run(port string) {
	log := utils.Logger()

	corsOpts := cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		MaxAge:         300,
	}

	r := chi.NewRouter()
	r.Use(chimiddleware.NoCache)
	r.Use(chimiddleware.SetHeader("Content-Type", "application/json"))
	r.Use(cors.Handler(corsOpts))

	// Create auth middleware
	authMiddleware := middleware.NewAuthMiddleware(s.service)

	r.Route("/api/core-back-end/v1", func(router chi.Router) {
		// Apply authentication middleware to all routes
		// Public endpoints will be skipped inside the middleware
		router.Use(authMiddleware.Authenticate)
		router.Mount("/", v1.Handler(s.handlers))
	})

	s.httpServer = &http.Server{
		Handler:           r,
		Addr:              port,
		ReadHeaderTimeout: 60 * time.Second,
	}

	log.WithField("port", port).Info("Server started")

	if err := s.httpServer.ListenAndServe(); err != nil {
		log.WithError(err).Fatal("failed to start server")
	}
}
