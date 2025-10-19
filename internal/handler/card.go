package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/google/uuid"
	openapi_types "github.com/oapi-codegen/runtime/types"
	v1 "github.com/tasks-control/core-back-end/api/v1"
	"github.com/tasks-control/core-back-end/internal/middleware"
	"github.com/tasks-control/core-back-end/internal/service"
	"github.com/tasks-control/core-back-end/pkg/utils"
)

// PostCards creates a new card
func (h *Handler) PostCards(w http.ResponseWriter, r *http.Request, params v1.PostCardsParams) {
	// Get authenticated user
	userID := middleware.MustGetUserIDFromContext(r.Context())

	// Parse request
	var req v1.CreateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Title == "" {
		utils.RespondError(w, http.StatusBadRequest, "Title is required")
		return
	}

	// Get list ID from query params
	listID := params.IdList

	// Convert position from float32 to float64 if provided
	// If position is 0, treat as nil to auto-calculate
	var positionPtr *float64
	if req.Position != nil && *req.Position != 0 {
		position := float64(*req.Position)
		positionPtr = &position
	}

	// Create card
	card, err := h.Service.CreateCard(r.Context(), service.CreateCardRequest{
		Title:       req.Title,
		Description: req.Description,
		IDList:      listID,
		Position:    positionPtr,
		MemberID:    userID,
	})
	if err != nil {
		if errors.Is(err, service.ErrNotBoardMember) {
			utils.RespondError(w, http.StatusForbidden, "You are not a member of this board")
			return
		}
		if errors.Is(err, service.ErrListNotFound) {
			utils.RespondError(w, http.StatusNotFound, "List not found")
			return
		}
		utils.Logger().WithError(err).Error("Failed to create card")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Convert to API response
	response := cardToAPIResponse(card)
	utils.RespondJSON(w, http.StatusCreated, response)
}

// GetCardsIdCard retrieves a card by ID
func (h *Handler) GetCardsIdCard(w http.ResponseWriter, r *http.Request, idCard openapi_types.UUID) {
	// Get authenticated user
	userID := middleware.MustGetUserIDFromContext(r.Context())

	// Get card
	card, err := h.Service.GetCard(r.Context(), idCard, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotBoardMember) {
			utils.RespondError(w, http.StatusForbidden, "You are not a member of this board")
			return
		}
		if errors.Is(err, service.ErrCardNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Card not found")
			return
		}
		utils.Logger().WithError(err).Error("Failed to get card")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Convert to API response
	response := cardToAPIResponse(card)
	utils.RespondJSON(w, http.StatusOK, response)
}

// PutCardsIdCard updates a card
func (h *Handler) PutCardsIdCard(w http.ResponseWriter, r *http.Request, idCard openapi_types.UUID) {
	// Get authenticated user
	userID := middleware.MustGetUserIDFromContext(r.Context())

	// Parse request
	var req v1.UpdateCardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Convert position from float32 to float64
	var positionPtr *float64
	if req.Position != nil {
		position := float64(*req.Position)
		positionPtr = &position
	}

	// Convert IDList from OpenAPI UUID to uuid.UUID
	var idListPtr *openapi_types.UUID
	if req.IdList != nil {
		idListPtr = req.IdList
	}

	// Update card
	card, err := h.Service.UpdateCard(r.Context(), idCard, userID, service.UpdateCardRequest{
		Title:       req.Title,
		Description: req.Description,
		IDList:      (*uuid.UUID)(idListPtr),
		Position:    positionPtr,
		Archived:    req.Archived,
	})
	if err != nil {
		if errors.Is(err, service.ErrNotBoardMember) {
			utils.RespondError(w, http.StatusForbidden, "You are not a member of this board")
			return
		}
		if errors.Is(err, service.ErrCardNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Card not found")
			return
		}
		if errors.Is(err, service.ErrListNotFound) {
			utils.RespondError(w, http.StatusNotFound, "List not found")
			return
		}
		utils.Logger().WithError(err).Error("Failed to update card")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Convert to API response
	response := cardToAPIResponse(card)
	utils.RespondJSON(w, http.StatusOK, response)
}

// DeleteCardsIdCard deletes a card
func (h *Handler) DeleteCardsIdCard(w http.ResponseWriter, r *http.Request, idCard openapi_types.UUID) {
	// Get authenticated user
	userID := middleware.MustGetUserIDFromContext(r.Context())

	// Delete card
	err := h.Service.DeleteCard(r.Context(), idCard, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotBoardMember) {
			utils.RespondError(w, http.StatusForbidden, "You are not a member of this board")
			return
		}
		if errors.Is(err, service.ErrCardNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Card not found")
			return
		}
		utils.Logger().WithError(err).Error("Failed to delete card")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
}
