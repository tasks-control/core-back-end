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

// GetBoards retrieves all boards for the authenticated user
func (h *Handler) GetBoards(w http.ResponseWriter, r *http.Request, params v1.GetBoardsParams) {
	// Get authenticated user
	userID := middleware.MustGetUserIDFromContext(r.Context())

	// Get query parameters
	starredOnly := false
	if params.Starred != nil {
		starredOnly = *params.Starred
	}

	limit := 20
	if params.Limit != nil {
		limit = *params.Limit
	}

	offset := 0
	if params.Offset != nil {
		offset = *params.Offset
	}

	// Get boards
	boards, total, err := h.Service.GetBoardsByMember(r.Context(), userID, starredOnly, limit, offset)
	if err != nil {
		utils.Logger().WithError(err).Error("Failed to get boards")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Convert to API response
	boardSummaries := make([]v1.BoardSummary, 0, len(boards))
	for _, board := range boards {
		boardSummaries = append(boardSummaries, boardToAPISummary(board))
	}

	response := struct {
		Boards []v1.BoardSummary `json:"boards"`
		Total  int               `json:"total"`
		Limit  int               `json:"limit"`
		Offset int               `json:"offset"`
	}{
		Boards: boardSummaries,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// PostBoards creates a new board
func (h *Handler) PostBoards(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user
	userID := middleware.MustGetUserIDFromContext(r.Context())

	// Parse request
	var req v1.CreateBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Name == "" {
		utils.RespondError(w, http.StatusBadRequest, "Name is required")
		return
	}
	if req.NameBoardUnique == "" {
		utils.RespondError(w, http.StatusBadRequest, "Name board unique is required")
		return
	}
	if req.Password == "" {
		utils.RespondError(w, http.StatusBadRequest, "Password is required")
		return
	}

	// Create board
	board, err := h.Service.CreateBoard(r.Context(), service.CreateBoardRequest{
		Name:            req.Name,
		NameBoardUnique: req.NameBoardUnique,
		Description:     req.Description,
		Password:        req.Password,
		CreatorID:       userID,
	})
	if err != nil {
		if errors.Is(err, service.ErrBoardAlreadyExists) {
			utils.RespondError(w, http.StatusConflict, "Board with this unique name already exists")
			return
		}
		if errors.Is(err, service.ErrInvalidBoardUniqueName) {
			utils.RespondError(w, http.StatusBadRequest, "Board unique name must contain only lowercase letters, numbers, and hyphens")
			return
		}
		utils.Logger().WithError(err).Error("Failed to create board")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Convert to API response
	response := boardToAPIResponse(board, false)
	utils.RespondJSON(w, http.StatusCreated, response)
}

// GetBoardsIdBoard retrieves a board with details
func (h *Handler) GetBoardsIdBoard(w http.ResponseWriter, r *http.Request, idBoard openapi_types.UUID) {
	// Get authenticated user
	userID := middleware.MustGetUserIDFromContext(r.Context())

	// Get board with details
	boardWithDetails, err := h.Service.GetBoardWithDetails(r.Context(), idBoard, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotBoardMember) {
			utils.RespondError(w, http.StatusForbidden, "You are not a member of this board")
			return
		}
		if errors.Is(err, service.ErrBoardNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Board not found")
			return
		}
		utils.Logger().WithError(err).Error("Failed to get board")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Convert to API response
	boardResp := boardToAPIResponse(boardWithDetails.Board, boardWithDetails.Board.Starred != nil && *boardWithDetails.Board.Starred)

	lists := make([]v1.List, 0, len(boardWithDetails.Lists))
	for _, list := range boardWithDetails.Lists {
		lists = append(lists, listToAPIResponse(list))
	}

	members := make([]v1.Member, 0, len(boardWithDetails.Members))
	for _, member := range boardWithDetails.Members {
		members = append(members, memberToAPIResponse(member))
	}

	response := struct {
		v1.Board
		Lists   []v1.List   `json:"lists"`
		Members []v1.Member `json:"members"`
	}{
		Board:   boardResp,
		Lists:   lists,
		Members: members,
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// PutBoardsIdBoard updates a board
func (h *Handler) PutBoardsIdBoard(w http.ResponseWriter, r *http.Request, idBoard openapi_types.UUID) {
	// Get authenticated user
	userID := middleware.MustGetUserIDFromContext(r.Context())

	// Parse request
	var req v1.UpdateBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Update board
	board, err := h.Service.UpdateBoard(r.Context(), idBoard, userID, service.UpdateBoardRequest{
		Name:            req.Name,
		NameBoardUnique: req.NameBoardUnique,
		Description:     req.Description,
		Password:        req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrNotBoardMember) {
			utils.RespondError(w, http.StatusForbidden, "You are not a member of this board")
			return
		}
		if errors.Is(err, service.ErrNotBoardOwner) {
			utils.RespondError(w, http.StatusForbidden, "Only board owners can update the board")
			return
		}
		if errors.Is(err, service.ErrBoardNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Board not found")
			return
		}
		if errors.Is(err, service.ErrBoardAlreadyExists) {
			utils.RespondError(w, http.StatusConflict, "Board with this unique name already exists")
			return
		}
		if errors.Is(err, service.ErrInvalidBoardUniqueName) {
			utils.RespondError(w, http.StatusBadRequest, "Board unique name must contain only lowercase letters, numbers, and hyphens")
			return
		}
		utils.Logger().WithError(err).Error("Failed to update board")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Convert to API response
	response := boardToAPIResponse(board, false)
	utils.RespondJSON(w, http.StatusOK, response)
}

// DeleteBoardsIdBoard deletes a board
func (h *Handler) DeleteBoardsIdBoard(w http.ResponseWriter, r *http.Request, idBoard openapi_types.UUID) {
	// Get authenticated user
	userID := middleware.MustGetUserIDFromContext(r.Context())

	// Delete board
	err := h.Service.DeleteBoard(r.Context(), idBoard, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotBoardMember) {
			utils.RespondError(w, http.StatusForbidden, "You are not a member of this board")
			return
		}
		if errors.Is(err, service.ErrNotBoardOwner) {
			utils.RespondError(w, http.StatusForbidden, "Only board owners can delete the board")
			return
		}
		if errors.Is(err, service.ErrBoardNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Board not found")
			return
		}
		utils.Logger().WithError(err).Error("Failed to delete board")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// PostMembersBoardsNameBoardUniqueJoin allows a user to join a board
func (h *Handler) PostMembersBoardsNameBoardUniqueJoin(w http.ResponseWriter, r *http.Request, nameBoardUnique string) {
	// Get authenticated user
	userID := middleware.MustGetUserIDFromContext(r.Context())

	// Parse request
	var req v1.JoinBoardRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Password == "" {
		utils.RespondError(w, http.StatusBadRequest, "Password is required")
		return
	}

	// Join board
	board, err := h.Service.JoinBoard(r.Context(), nameBoardUnique, req.Password, userID)
	if err != nil {
		if errors.Is(err, service.ErrBoardNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Board not found")
			return
		}
		if errors.Is(err, service.ErrInvalidBoardPassword) {
			utils.RespondError(w, http.StatusForbidden, "Invalid board password")
			return
		}
		if errors.Is(err, service.ErrAlreadyBoardMember) {
			utils.RespondError(w, http.StatusConflict, "Already a member of this board")
			return
		}
		utils.Logger().WithError(err).Error("Failed to join board")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Convert to API response
	boardResp := boardToAPIResponse(board, false)
	response := struct {
		Board v1.Board `json:"board"`
	}{
		Board: boardResp,
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// DeleteBoardsIdBoardMembersIdMember removes a member from a board
func (h *Handler) DeleteBoardsIdBoardMembersIdMember(w http.ResponseWriter, r *http.Request, idBoard openapi_types.UUID, idMember openapi_types.UUID) {
	// Get authenticated user
	userID := middleware.MustGetUserIDFromContext(r.Context())

	// Remove member
	err := h.Service.RemoveBoardMember(r.Context(), idBoard, idMember, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotBoardMember) {
			utils.RespondError(w, http.StatusForbidden, "You are not a member of this board")
			return
		}
		if errors.Is(err, service.ErrNotBoardOwner) {
			utils.RespondError(w, http.StatusForbidden, "Only board owners can remove members")
			return
		}
		if errors.Is(err, service.ErrCannotRemoveOwner) {
			utils.RespondError(w, http.StatusForbidden, "Cannot remove board owner")
			return
		}
		utils.Logger().WithError(err).Error("Failed to remove board member")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusOK)
}

// PostMembersBoardsIdBoardStar stars a board
func (h *Handler) PostMembersBoardsIdBoardStar(w http.ResponseWriter, r *http.Request, idBoard openapi_types.UUID) {
	// Get authenticated user
	userID := middleware.MustGetUserIDFromContext(r.Context())

	// Star board
	err := h.Service.StarBoard(r.Context(), idBoard, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotBoardMember) {
			utils.RespondError(w, http.StatusForbidden, "You are not a member of this board")
			return
		}
		if errors.Is(err, service.ErrBoardNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Board not found")
			return
		}
		utils.Logger().WithError(err).Error("Failed to star board")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	starred := true
	response := struct {
		Message string `json:"message"`
		Starred bool   `json:"starred"`
	}{
		Message: "Board starred successfully",
		Starred: starred,
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// DeleteMembersBoardsIdBoardStar unstars a board
func (h *Handler) DeleteMembersBoardsIdBoardStar(w http.ResponseWriter, r *http.Request, idBoard openapi_types.UUID) {
	// Get authenticated user
	userID := middleware.MustGetUserIDFromContext(r.Context())

	// Unstar board
	err := h.Service.UnstarBoard(r.Context(), idBoard, userID)
	if err != nil {
		if errors.Is(err, service.ErrNotBoardMember) {
			utils.RespondError(w, http.StatusForbidden, "You are not a member of this board")
			return
		}
		if errors.Is(err, service.ErrBoardNotFound) {
			utils.RespondError(w, http.StatusNotFound, "Board not found")
			return
		}
		utils.Logger().WithError(err).Error("Failed to unstar board")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	starred := false
	response := struct {
		Message string `json:"message"`
		Starred bool   `json:"starred"`
	}{
		Message: "Board unstarred successfully",
		Starred: starred,
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// Helper function to convert internal Board model to API response
func boardToAPIResponse(board *models.Board, starred bool) v1.Board {
	id := openapi_types.UUID(board.ID)
	idCreator := openapi_types.UUID(board.IDMemberCreator)

	return v1.Board{
		Id:              &id,
		Name:            &board.Name,
		NameBoardUnique: &board.NameBoardUnique,
		Description:     board.Description,
		IdMemberCreator: &idCreator,
		Starred:         &starred,
		CreatedAt:       &board.CreatedAt,
		UpdatedAt:       &board.UpdatedAt,
	}
}

// Helper function to convert internal Board model to API summary
func boardToAPISummary(board *models.Board) v1.BoardSummary {
	id := openapi_types.UUID(board.ID)
	starred := false
	if board.Starred != nil {
		starred = *board.Starred
	}

	summary := v1.BoardSummary{
		Id:              &id,
		Name:            &board.Name,
		NameBoardUnique: &board.NameBoardUnique,
		Description:     board.Description,
		Starred:         &starred,
		UpdatedAt:       &board.UpdatedAt,
	}

	if board.MemberCount != nil {
		summary.MemberCount = board.MemberCount
	}

	return summary
}

// Helper function to convert internal List model to API response
func listToAPIResponse(list *models.List) v1.List {
	id := openapi_types.UUID(list.ID)
	idBoard := openapi_types.UUID(list.IDBoard)
	position := float32(list.Position)

	return v1.List{
		Id:        &id,
		Name:      &list.Name,
		IdBoard:   &idBoard,
		Position:  &position,
		Archived:  &list.Archived,
		CreatedAt: &list.CreatedAt,
		UpdatedAt: &list.UpdatedAt,
	}
}
