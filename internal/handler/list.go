package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
	v1 "github.com/tasks-control/core-back-end/api/v1"
	"github.com/tasks-control/core-back-end/internal/middleware"
	"github.com/tasks-control/core-back-end/internal/models"
	"github.com/tasks-control/core-back-end/internal/service"
	"github.com/tasks-control/core-back-end/pkg/utils"
)

// PostLists creates a new list
func (h *Handler) PostLists(w http.ResponseWriter, r *http.Request, params v1.PostListsParams) {
	// Get authenticated user
	userID := middleware.MustGetUserIDFromContext(r.Context())

	// Parse request
	var req v1.CreateListRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Name == "" {
		utils.RespondError(w, http.StatusBadRequest, "Name is required")
		return
	}

	// Get board ID from query params
	boardID := params.IdBoard

	// Convert position from float32 to float64
	// If position is 0, treat as nil to auto-calculate
	var positionPtr *float64
	if req.Position != 0 {
		position := float64(req.Position)
		positionPtr = &position
	}

	// Create list
	list, err := h.Service.CreateList(r.Context(), service.CreateListRequest{
		Name:     req.Name,
		IDBoard:  boardID,
		Position: positionPtr,
		MemberID: userID,
	})
	if err != nil {
		if errors.Is(err, service.ErrNotBoardMember) {
			utils.RespondError(w, http.StatusForbidden, "You are not a member of this board")
			return
		}
		if errors.Is(err, service.ErrBoardNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Board not found")
			return
		}
		utils.Logger().WithError(err).Error("Failed to create list")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Convert to API response
	response := listToAPIResponse(list)
	utils.RespondJSON(w, http.StatusCreated, response)
}

// GetListsIdList retrieves a list with its cards
func (h *Handler) GetListsIdList(w http.ResponseWriter, r *http.Request, idList openapi_types.UUID) {
	// Get authenticated user
	userID := middleware.MustGetUserIDFromContext(r.Context())

	// Get list with cards
	listWithCards, err := h.Service.GetListWithCards(r.Context(), idList, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotBoardMember) {
			utils.RespondError(w, http.StatusForbidden, "You are not a member of this board")
			return
		}
		if errors.Is(err, service.ErrListNotFound) {
			utils.RespondError(w, http.StatusNotFound, "List not found")
			return
		}
		utils.Logger().WithError(err).Error("Failed to get list")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Convert to API response
	listResp := listToAPIResponse(listWithCards.List)

	cards := make([]v1.Card, 0, len(listWithCards.Cards))
	for _, card := range listWithCards.Cards {
		cards = append(cards, cardToAPIResponse(card))
	}

	response := struct {
		v1.List
		Cards []v1.Card `json:"cards"`
	}{
		List:  listResp,
		Cards: cards,
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// PutListsIdList updates a list
func (h *Handler) PutListsIdList(w http.ResponseWriter, r *http.Request, idList openapi_types.UUID) {
	// Get authenticated user
	userID := middleware.MustGetUserIDFromContext(r.Context())

	// Parse request
	var req v1.UpdateListRequest
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

	// Update list
	list, err := h.Service.UpdateList(r.Context(), idList, userID, service.UpdateListRequest{
		Name:     req.Name,
		Position: positionPtr,
		Archived: req.Archived,
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
		utils.Logger().WithError(err).Error("Failed to update list")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Convert to API response
	response := listToAPIResponse(list)
	utils.RespondJSON(w, http.StatusOK, response)
}

// DeleteListsIdList deletes a list
func (h *Handler) DeleteListsIdList(w http.ResponseWriter, r *http.Request, idList openapi_types.UUID) {
	// Get authenticated user
	userID := middleware.MustGetUserIDFromContext(r.Context())

	// Delete list
	err := h.Service.DeleteList(r.Context(), idList, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotBoardMember) {
			utils.RespondError(w, http.StatusForbidden, "You are not a member of this board")
			return
		}
		if errors.Is(err, service.ErrListNotFound) {
			utils.RespondError(w, http.StatusNotFound, "List not found")
			return
		}
		utils.Logger().WithError(err).Error("Failed to delete list")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Helper function to convert internal Card model to API response
func cardToAPIResponse(card *models.Card) v1.Card {
	id := openapi_types.UUID(card.ID)
	idList := openapi_types.UUID(card.IDList)
	createdBy := openapi_types.UUID(card.CreatedBy)
	position := float32(card.Position)

	return v1.Card{
		Id:          &id,
		Title:       &card.Title,
		Description: card.Description,
		IdList:      &idList,
		Position:    &position,
		Archived:    &card.Archived,
		CreatedBy:   &createdBy,
		CreatedAt:   &card.CreatedAt,
		UpdatedAt:   &card.UpdatedAt,
	}
}
