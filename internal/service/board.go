package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/google/uuid"
	"github.com/tasks-control/core-back-end/internal/models"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrBoardNotFound          = errors.New("board not found")
	ErrBoardAlreadyExists     = errors.New("board with this unique name already exists")
	ErrNotBoardMember         = errors.New("you are not a member of this board")
	ErrNotBoardOwner          = errors.New("you are not the owner of this board")
	ErrInvalidBoardPassword   = errors.New("invalid board password")
	ErrAlreadyBoardMember     = errors.New("already a member of this board")
	ErrCannotRemoveOwner      = errors.New("cannot remove board owner")
	ErrInvalidBoardUniqueName = errors.New("board unique name must contain only lowercase letters, numbers, and hyphens")
)

// CreateBoardRequest represents the data needed to create a new board
type CreateBoardRequest struct {
	Name            string
	NameBoardUnique string
	Description     *string
	Password        string
	CreatorID       uuid.UUID
}

// UpdateBoardRequest represents the data needed to update a board
type UpdateBoardRequest struct {
	Name            *string
	NameBoardUnique *string
	Description     *string
	Password        *string
}

// BoardWithDetails represents a board with its lists and members
type BoardWithDetails struct {
	Board   *models.Board
	Lists   []*models.List
	Members []*models.Member
}

// CreateBoard creates a new board
func (s *Service) CreateBoard(ctx context.Context, req CreateBoardRequest) (*models.Board, error) {
	// Validate unique name format
	if !isValidBoardUniqueName(req.NameBoardUnique) {
		return nil, ErrInvalidBoardUniqueName
	}

	// Check if unique name is already taken
	existingBoard, err := s.Repo.GetBoardByUniqueName(ctx, req.NameBoardUnique)
	if err != nil {
		return nil, fmt.Errorf("failed to check board existence: %w", err)
	}
	if existingBoard != nil {
		return nil, ErrBoardAlreadyExists
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	// Create board
	now := time.Now()
	board := &models.Board{
		ID:              uuid.New(),
		Name:            req.Name,
		NameBoardUnique: req.NameBoardUnique,
		Description:     req.Description,
		PasswordHash:    string(hashedPassword),
		IDMemberCreator: req.CreatorID,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	err = s.Repo.CreateBoard(ctx, board)
	if err != nil {
		return nil, fmt.Errorf("failed to create board: %w", err)
	}

	return board, nil
}

// GetBoardsByMember retrieves all boards a member belongs to
func (s *Service) GetBoardsByMember(ctx context.Context, memberID uuid.UUID, starredOnly bool, limit, offset int) ([]*models.Board, int, error) {
	// Validate pagination parameters
	if limit <= 0 || limit > 100 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}

	boards, total, err := s.Repo.GetBoardsByMemberID(ctx, memberID, starredOnly, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get boards: %w", err)
	}

	return boards, total, nil
}

// GetBoardWithDetails retrieves a board with its lists and members
func (s *Service) GetBoardWithDetails(ctx context.Context, boardID, memberID uuid.UUID) (*BoardWithDetails, error) {
	// Check if user is a member of the board
	boardMember, err := s.Repo.GetBoardMember(ctx, boardID, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to check board membership: %w", err)
	}
	if boardMember == nil {
		return nil, ErrNotBoardMember
	}

	// Get board
	board, err := s.Repo.GetBoardByID(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get board: %w", err)
	}
	if board == nil {
		return nil, ErrBoardNotFound
	}

	// Check if board is starred
	starred, err := s.Repo.IsStarred(ctx, boardID, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to check if board is starred: %w", err)
	}
	board.Starred = &starred

	// Get lists
	lists, err := s.Repo.GetBoardLists(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get board lists: %w", err)
	}

	// Get members
	members, err := s.Repo.GetBoardMembers(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get board members: %w", err)
	}

	return &BoardWithDetails{
		Board:   board,
		Lists:   lists,
		Members: members,
	}, nil
}

// UpdateBoard updates a board
func (s *Service) UpdateBoard(ctx context.Context, boardID, memberID uuid.UUID, req UpdateBoardRequest) (*models.Board, error) {
	// Check if user is a member of the board
	boardMember, err := s.Repo.GetBoardMember(ctx, boardID, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to check board membership: %w", err)
	}
	if boardMember == nil {
		return nil, ErrNotBoardMember
	}

	// Only owners can update board
	if boardMember.Role != models.BoardRoleOwner {
		return nil, ErrNotBoardOwner
	}

	// Get current board
	board, err := s.Repo.GetBoardByID(ctx, boardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get board: %w", err)
	}
	if board == nil {
		return nil, ErrBoardNotFound
	}

	// Update fields
	if req.Name != nil {
		board.Name = *req.Name
	}

	if req.NameBoardUnique != nil {
		// Validate unique name format
		if !isValidBoardUniqueName(*req.NameBoardUnique) {
			return nil, ErrInvalidBoardUniqueName
		}

		// Check if new unique name is already taken (if changed)
		if *req.NameBoardUnique != board.NameBoardUnique {
			existingBoard, err := s.Repo.GetBoardByUniqueName(ctx, *req.NameBoardUnique)
			if err != nil {
				return nil, fmt.Errorf("failed to check board existence: %w", err)
			}
			if existingBoard != nil {
				return nil, ErrBoardAlreadyExists
			}
		}

		board.NameBoardUnique = *req.NameBoardUnique
	}

	if req.Description != nil {
		board.Description = req.Description
	}

	if req.Password != nil {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, fmt.Errorf("failed to hash password: %w", err)
		}
		board.PasswordHash = string(hashedPassword)
	}

	board.UpdatedAt = time.Now()

	err = s.Repo.UpdateBoard(ctx, board)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrBoardNotFound
		}
		return nil, fmt.Errorf("failed to update board: %w", err)
	}

	return board, nil
}

// DeleteBoard deletes a board
func (s *Service) DeleteBoard(ctx context.Context, boardID, memberID uuid.UUID) error {
	// Check if user is a member of the board
	boardMember, err := s.Repo.GetBoardMember(ctx, boardID, memberID)
	if err != nil {
		return fmt.Errorf("failed to check board membership: %w", err)
	}
	if boardMember == nil {
		return ErrNotBoardMember
	}

	// Only owners can delete board
	if boardMember.Role != models.BoardRoleOwner {
		return ErrNotBoardOwner
	}

	err = s.Repo.DeleteBoard(ctx, boardID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrBoardNotFound
		}
		return fmt.Errorf("failed to delete board: %w", err)
	}

	return nil
}

// JoinBoard allows a user to join a board with password
func (s *Service) JoinBoard(ctx context.Context, uniqueName string, password string, memberID uuid.UUID) (*models.Board, error) {
	// Get board by unique name
	board, err := s.Repo.GetBoardByUniqueName(ctx, uniqueName)
	if err != nil {
		return nil, fmt.Errorf("failed to get board: %w", err)
	}
	if board == nil {
		return nil, ErrBoardNotFound
	}

	// Check if user is already a member
	existingMember, err := s.Repo.GetBoardMember(ctx, board.ID, memberID)
	if err != nil {
		return nil, fmt.Errorf("failed to check board membership: %w", err)
	}
	if existingMember != nil {
		return nil, ErrAlreadyBoardMember
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(board.PasswordHash), []byte(password))
	if err != nil {
		return nil, ErrInvalidBoardPassword
	}

	// Add member to board
	boardMember := &models.BoardMember{
		ID:       uuid.New(),
		IDBoard:  board.ID,
		IDMember: memberID,
		Role:     models.BoardRoleMember,
		JoinedAt: time.Now(),
	}

	err = s.Repo.AddBoardMember(ctx, boardMember)
	if err != nil {
		return nil, fmt.Errorf("failed to add member to board: %w", err)
	}

	return board, nil
}

// RemoveBoardMember removes a member from a board
func (s *Service) RemoveBoardMember(ctx context.Context, boardID, targetMemberID, requestingMemberID uuid.UUID) error {
	// Check if requesting user is a member of the board
	requestingBoardMember, err := s.Repo.GetBoardMember(ctx, boardID, requestingMemberID)
	if err != nil {
		return fmt.Errorf("failed to check board membership: %w", err)
	}
	if requestingBoardMember == nil {
		return ErrNotBoardMember
	}

	// Get target member
	targetBoardMember, err := s.Repo.GetBoardMember(ctx, boardID, targetMemberID)
	if err != nil {
		return fmt.Errorf("failed to check target board membership: %w", err)
	}
	if targetBoardMember == nil {
		return ErrNotBoardMember
	}

	// Check permissions
	// - Owner can remove anyone except themselves
	// - Member can only remove themselves
	if requestingMemberID != targetMemberID && requestingBoardMember.Role != models.BoardRoleOwner {
		return ErrNotBoardOwner
	}

	// Prevent removing the owner
	if targetBoardMember.Role == models.BoardRoleOwner && requestingMemberID == targetMemberID {
		return ErrCannotRemoveOwner
	}

	err = s.Repo.RemoveBoardMember(ctx, boardID, targetMemberID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrNotBoardMember
		}
		return fmt.Errorf("failed to remove member from board: %w", err)
	}

	return nil
}

// StarBoard stars a board for a member
func (s *Service) StarBoard(ctx context.Context, boardID, memberID uuid.UUID) error {
	// Check if user is a member of the board
	boardMember, err := s.Repo.GetBoardMember(ctx, boardID, memberID)
	if err != nil {
		return fmt.Errorf("failed to check board membership: %w", err)
	}
	if boardMember == nil {
		return ErrNotBoardMember
	}

	// Check if board exists
	board, err := s.Repo.GetBoardByID(ctx, boardID)
	if err != nil {
		return fmt.Errorf("failed to get board: %w", err)
	}
	if board == nil {
		return ErrBoardNotFound
	}

	err = s.Repo.StarBoard(ctx, boardID, memberID)
	if err != nil {
		return fmt.Errorf("failed to star board: %w", err)
	}

	return nil
}

// UnstarBoard unstars a board for a member
func (s *Service) UnstarBoard(ctx context.Context, boardID, memberID uuid.UUID) error {
	// Check if user is a member of the board
	boardMember, err := s.Repo.GetBoardMember(ctx, boardID, memberID)
	if err != nil {
		return fmt.Errorf("failed to check board membership: %w", err)
	}
	if boardMember == nil {
		return ErrNotBoardMember
	}

	// Check if board exists
	board, err := s.Repo.GetBoardByID(ctx, boardID)
	if err != nil {
		return fmt.Errorf("failed to get board: %w", err)
	}
	if board == nil {
		return ErrBoardNotFound
	}

	err = s.Repo.UnstarBoard(ctx, boardID, memberID)
	if err != nil {
		return fmt.Errorf("failed to unstar board: %w", err)
	}

	return nil
}

// isValidBoardUniqueName validates the board unique name format
func isValidBoardUniqueName(name string) bool {
	// Must contain only lowercase letters, numbers, and hyphens
	match, _ := regexp.MatchString(`^[a-z0-9-]+$`, name)
	return match && len(name) >= 3 && len(name) <= 50
}
