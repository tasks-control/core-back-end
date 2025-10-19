package handler

import (
	"net/http"

	"github.com/tasks-control/core-back-end/internal/service"

	openapi_types "github.com/oapi-codegen/runtime/types"
	v1 "github.com/tasks-control/core-back-end/api/v1"
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

func (h *Handler) PostCards(w http.ResponseWriter, r *http.Request, params v1.PostCardsParams) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) DeleteCardsIdCard(w http.ResponseWriter, r *http.Request, idCard openapi_types.UUID) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) GetCardsIdCard(w http.ResponseWriter, r *http.Request, idCard openapi_types.UUID) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}

func (h *Handler) PutCardsIdCard(w http.ResponseWriter, r *http.Request, idCard openapi_types.UUID) {
	//TODO implement me
	w.WriteHeader(http.StatusNotImplemented)
}
