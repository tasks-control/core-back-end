package server

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	v1 "github.com/tasks-control/core-back-end/api/v1"
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
}

func NewServer(handlers HTTPHandlers) *Server {
	return &Server{handlers: handlers}
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

	r.Route("/api/core-back-end/v1", func(router chi.Router) {
		//router.Use(middleware.GetAPIKeyMiddleware(s.gristService, s.cookieKey))
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
