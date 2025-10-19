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
	ErrListNotFound = errors.New("list not found")
)

// CreateListRequest represents the data needed to create a new list
type CreateListRequest struct {
	Name     string
	IDBoard  uuid.UUID
	Position *float64
	MemberID uuid.UUID
}

// UpdateListRequest represents the data needed to update a list
type UpdateListRequest struct {
	Name     *string
	Position *float64
	Archived *bool
}

// ListWithCards represents a list with its cards
type ListWithCards struct {
	List  *models.List
	Cards []*models.Card
}

// CreateList creates a new list in a board
func (s *Service) CreateList(ctx context.Context, req CreateListRequest) (*models.List, error) {
	// Check if user is a member of the board
	boardMember, err := s.Repo.GetBoardMember(ctx, req.IDBoard, req.MemberID)
	if err != nil {
		return nil, fmt.Errorf("failed to check board membership: %w", err)
	}
	if boardMember == nil {
		return nil, ErrNotBoardMember
	}

	// Check if board exists
	board, err := s.Repo.GetBoardByID(ctx, req.IDBoard)
	if err != nil {
		return nil, fmt.Errorf("failed to get board: %w", err)
	}
	if board == nil {
		return nil, ErrBoardNotFound
	}

	// Calculate position if not provided
	position := 0.0
	if req.Position != nil {
		position = *req.Position
	} else {
		// Auto-calculate position as max + 65536
		maxPos, err := s.Repo.GetMaxListPosition(ctx, req.IDBoard)
		if err != nil {
			return nil, fmt.Errorf("failed to get max list position: %w", err)
		}
		position = maxPos + 65536.0
	}

	// Create list
	now := time.Now()
	list := &models.List{
		ID:        uuid.New(),
		Name:      req.Name,
		IDBoard:   req.IDBoard,
		Position:  position,
		Archived:  false,
		CreatedAt: now,
		UpdatedAt: now,
	}

	err = s.Repo.CreateList(ctx, list)
	if err != nil {
		return nil, fmt.Errorf("failed to create list: %w", err)
	}

	return list, nil
}

// GetListWithCards retrieves a list with all its cards
func (s *Service) GetListWithCards(ctx context.Context, listID, memberID uuid.UUID) (*ListWithCards, error) {
	// Get list
	list, err := s.Repo.GetListByID(ctx, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list: %w", err)
	}
	if list == nil {
		return nil, ErrListNotFound
	}

	// Check if user is a member of the board
	boardMember, err := s.Repo.GetBoardMember(ctx, list.IDBoard, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to check board membership: %w", err)
	}
	if boardMember == nil {
		return nil, ErrNotBoardMember
	}

	// Get cards
	cards, err := s.Repo.GetListCards(ctx, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list cards: %w", err)
	}

	return &ListWithCards{
		List:  list,
		Cards: cards,
	}, nil
}

// UpdateList updates a list
func (s *Service) UpdateList(ctx context.Context, listID, memberID uuid.UUID, req UpdateListRequest) (*models.List, error) {
	// Get list
	list, err := s.Repo.GetListByID(ctx, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list: %w", err)
	}
	if list == nil {
		return nil, ErrListNotFound
	}

	// Check if user is a member of the board
	boardMember, err := s.Repo.GetBoardMember(ctx, list.IDBoard, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to check board membership: %w", err)
	}
	if boardMember == nil {
		return nil, ErrNotBoardMember
	}

	// Update fields
	if req.Name != nil {
		list.Name = *req.Name
	}

	if req.Position != nil {
		list.Position = *req.Position
	}

	if req.Archived != nil {
		list.Archived = *req.Archived
	}

	list.UpdatedAt = time.Now()

	// Save updated list
	err = s.Repo.UpdateList(ctx, list)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrListNotFound
		}
		return nil, fmt.Errorf("failed to update list: %w", err)
	}

	return list, nil
}

// DeleteList deletes a list
func (s *Service) DeleteList(ctx context.Context, listID, memberID uuid.UUID) error {
	// Get list
	list, err := s.Repo.GetListByID(ctx, listID)
	if err != nil {
		return fmt.Errorf("failed to get list: %w", err)
	}
	if list == nil {
		return ErrListNotFound
	}

	// Check if user is a member of the board
	boardMember, err := s.Repo.GetBoardMember(ctx, list.IDBoard, memberID)
	if err != nil {
		return fmt.Errorf("failed to check board membership: %w", err)
	}
	if boardMember == nil {
		return ErrNotBoardMember
	}

	// Delete list
	err = s.Repo.DeleteList(ctx, listID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrListNotFound
		}
		return fmt.Errorf("failed to delete list: %w", err)
	}

	return nil
}
