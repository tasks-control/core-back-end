package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/tasks-control/core-back-end/internal/models"
	"github.com/tasks-control/core-back-end/internal/service"
	"github.com/tasks-control/core-back-end/pkg/utils"
)

// ContextKey is a custom type for context keys to avoid collisions
type ContextKey string

const (
	// UserContextKey is the key for storing user info in request context
	UserContextKey ContextKey = "user"
	// UserIDContextKey is the key for storing user ID in request context
	UserIDContextKey ContextKey = "user_id"
)

// AuthMiddleware creates a middleware that validates JWT tokens
type AuthMiddleware struct {
	service *service.Service
}

// NewAuthMiddleware creates a new auth middleware
func NewAuthMiddleware(svc *service.Service) *AuthMiddleware {
	return &AuthMiddleware{
		service: svc,
	}
}

// Authenticate is a middleware that validates JWT tokens and loads user info
// It skips authentication for public endpoints (those with security: [] in OpenAPI spec)
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// List of public endpoints that don't require authentication
		publicEndpoints := []string{
			"/auth/register",
			"/auth/login",
			"/auth/refresh",
			"/alive",
		}

		// Check if current path is a public endpoint
		for _, endpoint := range publicEndpoints {
			if strings.HasSuffix(r.URL.Path, endpoint) {
				next.ServeHTTP(w, r)
				return
			}
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.RespondError(w, http.StatusUnauthorized, "Authorization header required")
			return
		}

		// Check if it's a Bearer token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			utils.RespondError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		token := parts[1]

		// Validate token
		claims, err := m.service.JWTManager.ValidateAccessToken(token)
		if err != nil {
			if err == utils.ErrExpiredToken {
				utils.RespondError(w, http.StatusUnauthorized, "Token has expired")
				return
			}
			utils.RespondError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Load user from database to ensure they still exist and get fresh data
		user, err := m.service.GetMemberByID(r.Context(), claims.UserID)
		if err != nil {
			utils.Logger().WithError(err).Error("Failed to load user from database")
			utils.RespondError(w, http.StatusInternalServerError, "Internal server error")
			return
		}

		if user == nil {
			utils.RespondError(w, http.StatusUnauthorized, "User not found")
			return
		}

		// Add user info to context
		ctx := context.WithValue(r.Context(), UserContextKey, user)
		ctx = context.WithValue(ctx, UserIDContextKey, claims.UserID)

		// Call next handler with updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext extracts the user from the request context
func GetUserFromContext(ctx context.Context) (*models.Member, bool) {
	user, ok := ctx.Value(UserContextKey).(*models.Member)
	return user, ok
}

// GetUserIDFromContext extracts the user ID from the request context
func GetUserIDFromContext(ctx context.Context) (uuid.UUID, bool) {
	userID, ok := ctx.Value(UserIDContextKey).(uuid.UUID)
	return userID, ok
}

// MustGetUserFromContext extracts the user from context or panics
// Use only in handlers that are guaranteed to have authentication middleware
func MustGetUserFromContext(ctx context.Context) *models.Member {
	user, ok := GetUserFromContext(ctx)
	if !ok {
		panic("user not found in context - authentication middleware not applied?")
	}
	return user
}

// MustGetUserIDFromContext extracts the user ID from context or panics
// Use only in handlers that are guaranteed to have authentication middleware
func MustGetUserIDFromContext(ctx context.Context) uuid.UUID {
	userID, ok := GetUserIDFromContext(ctx)
	if !ok {
		panic("user ID not found in context - authentication middleware not applied?")
	}
	return userID
}
