package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tasks-control/core-back-end/internal/models"
)

var (
	ErrCardNotFound = errors.New("card not found")
)

// CreateCardRequest represents the data needed to create a new card
type CreateCardRequest struct {
	Title       string
	Description *string
	IDList      uuid.UUID
	Position    *float64
	MemberID    uuid.UUID
}

// UpdateCardRequest represents the data needed to update a card
type UpdateCardRequest struct {
	Title       *string
	Description *string
	IDList      *uuid.UUID // Allow moving to different list
	Position    *float64
	Archived    *bool
}

// CreateCard creates a new card in a list
func (s *Service) CreateCard(ctx context.Context, req CreateCardRequest) (*models.Card, error) {
	// Get board ID from list
	boardID, err := s.Repo.GetBoardIDByListID(ctx, req.IDList)
	if err != nil {
		return nil, fmt.Errorf("failed to get board ID: %w", err)
	}
	if boardID == uuid.Nil {
		return nil, ErrListNotFound
	}

	// Check if user is a member of the board
	boardMember, err := s.Repo.GetBoardMember(ctx, boardID, req.MemberID)
	if err != nil {
		return nil, fmt.Errorf("failed to check board membership: %w", err)
	}
	if boardMember == nil {
		return nil, ErrNotBoardMember
	}

	// Check if list exists
	list, err := s.Repo.GetListByID(ctx, req.IDList)
	if err != nil {
		return nil, fmt.Errorf("failed to get list: %w", err)
	}
	if list == nil {
		return nil, ErrListNotFound
	}

	// Calculate position if not provided
	position := 0.0
	if req.Position != nil {
		position = *req.Position
	} else {
		// Auto-calculate position as max + 65536
		maxPos, err := s.Repo.GetMaxCardPosition(ctx, req.IDList)
		if err != nil {
			return nil, fmt.Errorf("failed to get max card position: %w", err)
		}
		position = maxPos + 65536.0
	}

	// Create card
	now := time.Now()
	card := &models.Card{
		ID:          uuid.New(),
		Title:       req.Title,
		Description: req.Description,
		IDList:      req.IDList,
		Position:    position,
		Archived:    false,
		CreatedBy:   req.MemberID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	err = s.Repo.CreateCard(ctx, card)
	if err != nil {
		return nil, fmt.Errorf("failed to create card: %w", err)
	}

	return card, nil
}

// GetCard retrieves a card by ID
func (s *Service) GetCard(ctx context.Context, cardID, memberID uuid.UUID) (*models.Card, error) {
	// Get card
	card, err := s.Repo.GetCardByID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}
	if card == nil {
		return nil, ErrCardNotFound
	}

	// Get board ID from card
	boardID, err := s.Repo.GetBoardIDByCardID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get board ID: %w", err)
	}
	if boardID == uuid.Nil {
		return nil, ErrCardNotFound
	}

	// Check if user is a member of the board
	boardMember, err := s.Repo.GetBoardMember(ctx, boardID, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to check board membership: %w", err)
	}
	if boardMember == nil {
		return nil, ErrNotBoardMember
	}

	return card, nil
}

// UpdateCard updates a card
func (s *Service) UpdateCard(ctx context.Context, cardID, memberID uuid.UUID, req UpdateCardRequest) (*models.Card, error) {
	// Get card
	card, err := s.Repo.GetCardByID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}
	if card == nil {
		return nil, ErrCardNotFound
	}

	// Get board ID from card
	boardID, err := s.Repo.GetBoardIDByCardID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get board ID: %w", err)
	}
	if boardID == uuid.Nil {
		return nil, ErrCardNotFound
	}

	// Check if user is a member of the board
	boardMember, err := s.Repo.GetBoardMember(ctx, boardID, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to check board membership: %w", err)
	}
	if boardMember == nil {
		return nil, ErrNotBoardMember
	}

	// Update fields
	if req.Title != nil {
		card.Title = *req.Title
	}

	if req.Description != nil {
		card.Description = req.Description
	}

	// Handle moving card to different list
	if req.IDList != nil && *req.IDList != card.IDList {
		// Check if target list exists and is in the same board
		targetList, err := s.Repo.GetListByID(ctx, *req.IDList)
		if err != nil {
			return nil, fmt.Errorf("failed to get target list: %w", err)
		}
		if targetList == nil {
			return nil, ErrListNotFound
		}

		// Verify target list is in the same board
		if targetList.IDBoard != boardID {
			return nil, fmt.Errorf("cannot move card to a list in a different board")
		}

		card.IDList = *req.IDList

		// If position not specified when moving, auto-calculate for new list
		if req.Position == nil {
			maxPos, err := s.Repo.GetMaxCardPosition(ctx, *req.IDList)
			if err != nil {
				return nil, fmt.Errorf("failed to get max card position: %w", err)
			}
			card.Position = maxPos + 65536.0
		}
	}

	if req.Position != nil {
		card.Position = *req.Position
	}

	if req.Archived != nil {
		card.Archived = *req.Archived
	}

	card.UpdatedAt = time.Now()

	// Save updated card
	err = s.Repo.UpdateCard(ctx, card)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrCardNotFound
		}
		return nil, fmt.Errorf("failed to update card: %w", err)
	}

	return card, nil
}

// DeleteCard deletes a card
func (s *Service) DeleteCard(ctx context.Context, cardID, memberID uuid.UUID) error {
	// Get card
	card, err := s.Repo.GetCardByID(ctx, cardID)
	if err != nil {
		return fmt.Errorf("failed to get card: %w", err)
	}
	if card == nil {
		return ErrCardNotFound
	}

	// Get board ID from card
	boardID, err := s.Repo.GetBoardIDByCardID(ctx, cardID)
	if err != nil {
		return fmt.Errorf("failed to get board ID: %w", err)
	}
	if boardID == uuid.Nil {
		return ErrCardNotFound
	}

	// Check if user is a member of the board
	boardMember, err := s.Repo.GetBoardMember(ctx, boardID, memberID)
	if err != nil {
		return fmt.Errorf("failed to check board membership: %w", err)
	}
	if boardMember == nil {
		return ErrNotBoardMember
	}

	// Delete card
	err = s.Repo.DeleteCard(ctx, cardID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrCardNotFound
		}
		return fmt.Errorf("failed to delete card: %w", err)
	}

	return nil
}
