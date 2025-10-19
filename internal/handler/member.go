package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	v1 "github.com/tasks-control/core-back-end/api/v1"
	"github.com/tasks-control/core-back-end/internal/middleware"
	"github.com/tasks-control/core-back-end/internal/service"
	"github.com/tasks-control/core-back-end/pkg/utils"
)

// GetMembersMe retrieves the current authenticated user's information
func (h *Handler) GetMembersMe(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user
	user := middleware.MustGetUserFromContext(r.Context())

	// Convert to API response
	response := memberToAPIResponse(user)
	utils.RespondJSON(w, http.StatusOK, response)
}

// PutMembersMe updates the current authenticated user's profile
func (h *Handler) PutMembersMe(w http.ResponseWriter, r *http.Request) {
	// Get authenticated user
	userID := middleware.MustGetUserIDFromContext(r.Context())

	// Parse request
	var req v1.UpdateMemberRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate email format if provided
	if req.Email != nil {
		email := string(*req.Email)
		if email == "" {
			utils.RespondError(w, http.StatusBadRequest, "Email cannot be empty")
			return
		}
	}

	// Validate username if provided
	if req.Username != nil && *req.Username == "" {
		utils.RespondError(w, http.StatusBadRequest, "Username cannot be empty")
		return
	}

	// Validate password length if provided
	if req.Password != nil && *req.Password != "" && len(*req.Password) < 8 {
		utils.RespondError(w, http.StatusBadRequest, "Password must be at least 8 characters long")
		return
	}

	// Convert email from OpenAPI type to string pointer
	var emailPtr *string
	if req.Email != nil {
		email := string(*req.Email)
		emailPtr = &email
	}

	// Update member profile
	member, err := h.Service.UpdateMemberProfile(r.Context(), userID, service.UpdateMemberRequest{
		Email:    emailPtr,
		Username: req.Username,
		FullName: req.FullName,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserNotFound) {
			utils.RespondError(w, http.StatusNotFound, "User not found")
			return
		}
		if errors.Is(err, service.ErrEmailAlreadyTaken) {
			utils.RespondError(w, http.StatusConflict, "Email is already taken")
			return
		}
		if errors.Is(err, service.ErrUsernameAlreadyTaken) {
			utils.RespondError(w, http.StatusConflict, "Username is already taken")
			return
		}
		utils.Logger().WithError(err).Error("Failed to update member profile")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Convert to API response
	response := memberToAPIResponse(member)
	utils.RespondJSON(w, http.StatusOK, response)
}
