package handler

import (
	"net/http"

	"github.com/tasks-control/core-back-end/internal/service"
)

type Handler struct {
	Service *service.Service
}

func NewHandler(s *service.Service) *Handler {
	return &Handler{Service: s}
}

func (h *Handler) GetAlive(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) Readiness(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}
