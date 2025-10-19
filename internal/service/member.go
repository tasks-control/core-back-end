package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tasks-control/core-back-end/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrEmailAlreadyTaken    = errors.New("email is already taken")
	ErrUsernameAlreadyTaken = errors.New("username is already taken")
)

// UpdateMemberRequest represents the data needed to update a member profile
type UpdateMemberRequest struct {
	Email    *string
	Username *string
	FullName *string
	Password *string
}

// GetMemberProfile retrieves the member profile by ID
func (s *Service) GetMemberProfile(ctx context.Context, memberID uuid.UUID) (*models.Member, error) {
	member, err := s.Repo.GetMemberByID(ctx, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to get member profile: %w", err)
	}
	if member == nil {
		return nil, ErrUserNotFound
	}
	return member, nil
}

// UpdateMemberProfile updates a member's profile information
func (s *Service) UpdateMemberProfile(ctx context.Context, memberID uuid.UUID, req UpdateMemberRequest) (*models.Member, error) {
	// Get current member
	member, err := s.Repo.GetMemberByID(ctx, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to get member: %w", err)
	}
	if member == nil {
		return nil, ErrUserNotFound
	}

	// Update email if provided
	if req.Email != nil && *req.Email != member.Email {
		// Check if email is already taken
		existingMember, err := s.Repo.GetMemberByEmail(ctx, *req.Email)
		if err != nil {
			return nil, fmt.Errorf("failed to check email uniqueness: %w", err)
		}
		if existingMember != nil && existingMember.ID != memberID {
			return nil, ErrEmailAlreadyTaken
		}
		member.Email = *req.Email
	}

	// Update username if provided
	if req.Username != nil && *req.Username != member.Username {
		// Check if username is already taken
		existingMember, err := s.Repo.GetMemberByUsername(ctx, *req.Username)
		if err != nil {
			return nil, fmt.Errorf("failed to check username uniqueness: %w", err)
		}
		if existingMember != nil && existingMember.ID != memberID {
			return nil, ErrUsernameAlreadyTaken
		}
		member.Username = *req.Username
	}

	// Update full name if provided
	if req.FullName != nil {
		member.FullName = req.FullName
	}

	// Update password if provided
	if req.Password != nil && *req.Password != "" {
		// Validate password length
		if len(*req.Password) < 8 {
			return nil, fmt.Errorf("password must be at least 8 characters long")
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		member.PasswordHash = string(hashedPassword)
	}

	member.UpdatedAt = time.Now()

	// Save updated member
	err = s.Repo.UpdateMember(ctx, member)
	if err != nil {
		return nil, fmt.Errorf("failed to update member: %w", err)
	}

	return member, nil
}
