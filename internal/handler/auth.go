package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	openapi_types "github.com/oapi-codegen/runtime/types"
	v1 "github.com/tasks-control/core-back-end/api/v1"
	"github.com/tasks-control/core-back-end/internal/models"
	"github.com/tasks-control/core-back-end/internal/service"
	"github.com/tasks-control/core-back-end/pkg/utils"
)

// PostAuthRegister handles user registration
func (h *Handler) PostAuthRegister(w http.ResponseWriter, r *http.Request) {
	var req v1.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Email == "" {
		utils.RespondError(w, http.StatusBadRequest, "Email is required")
		return
	}
	if req.Password == "" {
		utils.RespondError(w, http.StatusBadRequest, "Password is required")
		return
	}
	if req.Username == "" {
		utils.RespondError(w, http.StatusBadRequest, "Username is required")
		return
	}

	// Validate password length
	if len(req.Password) < 8 {
		utils.RespondError(w, http.StatusBadRequest, "Password must be at least 8 characters long")
		return
	}

	// Register user
	member, err := h.Service.Register(r.Context(), service.RegisterRequest{
		Email:    string(req.Email),
		Password: req.Password,
		Username: req.Username,
		FullName: req.FullName,
	})
	if err != nil {
		if errors.Is(err, service.ErrUserAlreadyExists) {
			utils.RespondError(w, http.StatusConflict, "User already exists")
			return
		}
		utils.Logger().WithError(err).Error("Failed to register user")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Convert to response format
	response := memberToAPIResponse(member)
	utils.RespondJSON(w, http.StatusCreated, response)
}

// PostAuthLogin handles user login
func (h *Handler) PostAuthLogin(w http.ResponseWriter, r *http.Request) {
	var req v1.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Email == "" {
		utils.RespondError(w, http.StatusBadRequest, "Email is required")
		return
	}
	if req.Password == "" {
		utils.RespondError(w, http.StatusBadRequest, "Password is required")
		return
	}

	// Login user
	authResp, err := h.Service.Login(r.Context(), service.LoginRequest{
		Email:    string(req.Email),
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid credentials")
			return
		}
		utils.Logger().WithError(err).Error("Failed to login user")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Convert to response format
	response := v1.LoginResponse{
		AccessToken:  &authResp.AccessToken,
		RefreshToken: &authResp.RefreshToken,
		ExpiresIn:    &authResp.ExpiresIn,
		User:         memberToAPIResponsePtr(authResp.User),
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// PostAuthRefresh handles access token refresh
func (h *Handler) PostAuthRefresh(w http.ResponseWriter, r *http.Request) {
	var req v1.TokenRefreshRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.RefreshToken == "" {
		utils.RespondError(w, http.StatusBadRequest, "Refresh token is required")
		return
	}

	// Refresh access token
	authResp, err := h.Service.RefreshAccessToken(r.Context(), req.RefreshToken)
	if err != nil {
		if errors.Is(err, service.ErrInvalidRefreshToken) || errors.Is(err, utils.ErrInvalidToken) || errors.Is(err, utils.ErrExpiredToken) {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
			return
		}
		utils.Logger().WithError(err).Error("Failed to refresh token")
		utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	// Convert to response format
	response := v1.TokenRefreshResponse{
		AccessToken: &authResp.AccessToken,
		ExpiresIn:   &authResp.ExpiresIn,
	}

	utils.RespondJSON(w, http.StatusOK, response)
}

// Helper function to convert internal Member model to API response
func memberToAPIResponse(member *models.Member) v1.Member {
	email := openapi_types.Email(member.Email)
	id := openapi_types.UUID(member.ID)

	return v1.Member{
		Id:        &id,
		Email:     &email,
		Username:  &member.Username,
		FullName:  member.FullName,
		CreatedAt: &member.CreatedAt,
		UpdatedAt: &member.UpdatedAt,
	}
}

// Helper function to convert internal Member model to API response pointer
func memberToAPIResponsePtr(member *models.Member) *v1.Member {
	resp := memberToAPIResponse(member)
	return &resp
}
