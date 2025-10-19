package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/tasks-control/core-back-end/internal/models"
	"github.com/tasks-control/core-back-end/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// RegisterRequest represents the data needed to register a new user
type RegisterRequest struct {
	Email    string
	Password string
	Username string
	FullName *string
}

// LoginRequest represents the data needed to login
type LoginRequest struct {
	Email    string
	Password string
}

// AuthResponse represents the response after successful authentication
type AuthResponse struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	User         *models.Member
}

// Register creates a new user account
func (s *Service) Register(ctx context.Context, req RegisterRequest) (*models.Member, error) {
	// Check if user already exists by email
	existingUser, err := s.Repo.GetMemberByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Check if username is taken
	existingUser, err = s.Repo.GetMemberByUsername(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create new member
	now := time.Now()
	member := &models.Member{
		ID:           uuid.New(),
		Email:        req.Email,
		Username:     req.Username,
		FullName:     req.FullName,
		PasswordHash: string(hashedPassword),
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	// Save to a database
	err = s.Repo.CreateMember(ctx, member)
	if err != nil {
		return nil, err
	}

	return member, nil
}

// Login authenticates a user and returns tokens
func (s *Service) Login(ctx context.Context, req LoginRequest) (*AuthResponse, error) {
	// Get user by email
	member, err := s.Repo.GetMemberByEmail(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrInvalidCredentials
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(member.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Generate tokens
	accessTokenDuration := time.Duration(s.JWTConfig.AccessTokenDuration) * time.Second
	refreshTokenDuration := time.Duration(s.JWTConfig.RefreshTokenDuration) * time.Second

	accessToken, err := s.JWTManager.GenerateAccessToken(
		member.ID,
		member.Email,
		member.Username,
		accessTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.JWTManager.GenerateRefreshToken(member.ID, refreshTokenDuration)
	if err != nil {
		return nil, err
	}

	// Store refresh token in database
	tokenHash := utils.HashToken(refreshToken)
	refreshTokenModel := &models.RefreshToken{
		ID:        uuid.New(),
		IDMember:  member.ID,
		TokenHash: tokenHash,
		ExpiresAt: time.Now().Add(refreshTokenDuration),
		CreatedAt: time.Now(),
		Revoked:   false,
	}

	err = s.Repo.CreateRefreshToken(ctx, refreshTokenModel)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    s.JWTConfig.AccessTokenDuration,
		User:         member,
	}, nil
}

// RefreshAccessToken generates a new access token using a refresh token
func (s *Service) RefreshAccessToken(ctx context.Context, refreshToken string) (*AuthResponse, error) {
	// Validate refresh token
	userID, err := s.JWTManager.ValidateRefreshToken(refreshToken)
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}

	// Check if refresh token exists and is valid in a database
	tokenHash := utils.HashToken(refreshToken)
	storedToken, err := s.Repo.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}
	if storedToken == nil {
		return nil, ErrInvalidRefreshToken
	}

	// Get user
	member, err := s.GetMemberByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if member == nil {
		return nil, ErrUserNotFound
	}

	// Generate new access token
	accessTokenDuration := time.Duration(s.JWTConfig.AccessTokenDuration) * time.Second
	accessToken, err := s.JWTManager.GenerateAccessToken(
		member.ID,
		member.Email,
		member.Username,
		accessTokenDuration,
	)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		AccessToken: accessToken,
		ExpiresIn:   s.JWTConfig.AccessTokenDuration,
		User:        member,
	}, nil
}

// GetMemberByID retrieves a member by ID (useful for authentication middleware)
func (s *Service) GetMemberByID(ctx context.Context, id uuid.UUID) (*models.Member, error) {
	return s.Repo.GetMemberByID(ctx, id)
}
