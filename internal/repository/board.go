package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tasks-control/core-back-end/internal/models"
)

// CreateBoard inserts a new board into the database and adds the creator as owner
func (r *repository) CreateBoard(ctx context.Context, board *models.Board) error {
	tx, err := r.conn.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Insert board
	query := `
		INSERT INTO boards (id, name, name_board_unique, description, password_hash, id_member_creator, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err = tx.ExecContext(ctx, query,
		board.ID,
		board.Name,
		board.NameBoardUnique,
		board.Description,
		board.PasswordHash,
		board.IDMemberCreator,
		board.CreatedAt,
		board.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create board: %w", err)
	}

	// Add creator as owner in board_members
	boardMemberQuery := `
		INSERT INTO board_members (id, id_board, id_member, role, joined_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err = tx.ExecContext(ctx, boardMemberQuery,
		uuid.New(),
		board.ID,
		board.IDMemberCreator,
		models.BoardRoleOwner,
		board.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to add creator as board member: %w", err)
	}

	return tx.Commit()
}

// GetBoardByID retrieves a board by ID
func (r *repository) GetBoardByID(ctx context.Context, boardID uuid.UUID) (*models.Board, error) {
	var board models.Board
	query := `
		SELECT id, name, name_board_unique, description, password_hash, id_member_creator, created_at, updated_at
		FROM boards
		WHERE id = $1
	`
	err := r.conn.GetContext(ctx, &board, query, boardID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get board: %w", err)
	}
	return &board, nil
}

// GetBoardByUniqueName retrieves a board by its unique name
func (r *repository) GetBoardByUniqueName(ctx context.Context, uniqueName string) (*models.Board, error) {
	var board models.Board
	query := `
		SELECT id, name, name_board_unique, description, password_hash, id_member_creator, created_at, updated_at
		FROM boards
		WHERE name_board_unique = $1
	`
	err := r.conn.GetContext(ctx, &board, query, uniqueName)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get board by unique name: %w", err)
	}
	return &board, nil
}

// GetBoardsByMemberID retrieves all boards a member belongs to with pagination
func (r *repository) GetBoardsByMemberID(ctx context.Context, memberID uuid.UUID, starredOnly bool, limit, offset int) ([]*models.Board, int, error) {
	boards := []*models.Board{}

	// Build query based on filters
	query := `
		SELECT 
			b.id, 
			b.name, 
			b.name_board_unique, 
			b.description, 
			b.password_hash,
			b.id_member_creator, 
			b.created_at, 
			b.updated_at,
			EXISTS(SELECT 1 FROM starred_boards sb WHERE sb.id_board = b.id AND sb.id_member = $1) as starred,
			(SELECT COUNT(*) FROM board_members WHERE id_board = b.id) as member_count
		FROM boards b
		INNER JOIN board_members bm ON b.id = bm.id_board
	`

	countQuery := `
		SELECT COUNT(*)
		FROM boards b
		INNER JOIN board_members bm ON b.id = bm.id_board
	`

	whereClause := " WHERE bm.id_member = $1"

	if starredOnly {
		whereClause += " AND EXISTS(SELECT 1 FROM starred_boards sb WHERE sb.id_board = b.id AND sb.id_member = $1)"
	}

	query += whereClause + " ORDER BY b.updated_at DESC LIMIT $2 OFFSET $3"
	countQuery += whereClause

	// Get total count
	var total int
	err := r.conn.GetContext(ctx, &total, countQuery, memberID)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count boards: %w", err)
	}

	// Get boards
	err = r.conn.SelectContext(ctx, &boards, query, memberID, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get boards: %w", err)
	}

	return boards, total, nil
}

// UpdateBoard updates an existing board
func (r *repository) UpdateBoard(ctx context.Context, board *models.Board) error {
	query := `
		UPDATE boards
		SET name = $2, name_board_unique = $3, description = $4, password_hash = $5, updated_at = $6
		WHERE id = $1
	`
	result, err := r.conn.ExecContext(ctx, query,
		board.ID,
		board.Name,
		board.NameBoardUnique,
		board.Description,
		board.PasswordHash,
		board.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update board: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// DeleteBoard deletes a board (cascade will delete related data)
func (r *repository) DeleteBoard(ctx context.Context, boardID uuid.UUID) error {
	query := `DELETE FROM boards WHERE id = $1`
	result, err := r.conn.ExecContext(ctx, query, boardID)
	if err != nil {
		return fmt.Errorf("failed to delete board: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetBoardMember retrieves a specific board member relationship
func (r *repository) GetBoardMember(ctx context.Context, boardID, memberID uuid.UUID) (*models.BoardMember, error) {
	var boardMember models.BoardMember
	query := `
		SELECT id, id_board, id_member, role, joined_at
		FROM board_members
		WHERE id_board = $1 AND id_member = $2
	`
	err := r.conn.GetContext(ctx, &boardMember, query, boardID, memberID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get board member: %w", err)
	}
	return &boardMember, nil
}

// AddBoardMember adds a member to a board
func (r *repository) AddBoardMember(ctx context.Context, boardMember *models.BoardMember) error {
	query := `
		INSERT INTO board_members (id, id_board, id_member, role, joined_at)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.conn.ExecContext(ctx, query,
		boardMember.ID,
		boardMember.IDBoard,
		boardMember.IDMember,
		boardMember.Role,
		boardMember.JoinedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to add board member: %w", err)
	}
	return nil
}

// RemoveBoardMember removes a member from a board
func (r *repository) RemoveBoardMember(ctx context.Context, boardID, memberID uuid.UUID) error {
	query := `DELETE FROM board_members WHERE id_board = $1 AND id_member = $2`
	result, err := r.conn.ExecContext(ctx, query, boardID, memberID)
	if err != nil {
		return fmt.Errorf("failed to remove board member: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

// GetBoardMembers retrieves all members of a board
func (r *repository) GetBoardMembers(ctx context.Context, boardID uuid.UUID) ([]*models.Member, error) {
	members := []*models.Member{}
	query := `
		SELECT m.id, m.email, m.username, m.full_name, m.password_hash, m.created_at, m.updated_at
		FROM members m
		INNER JOIN board_members bm ON m.id = bm.id_member
		WHERE bm.id_board = $1
		ORDER BY bm.joined_at ASC
	`
	err := r.conn.SelectContext(ctx, &members, query, boardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get board members: %w", err)
	}
	return members, nil
}

// StarBoard adds a star to a board for a member
func (r *repository) StarBoard(ctx context.Context, boardID, memberID uuid.UUID) error {
	query := `
		INSERT INTO starred_boards (id, id_board, id_member, starred_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (id_board, id_member) DO NOTHING
	`
	_, err := r.conn.ExecContext(ctx, query, uuid.New(), boardID, memberID, time.Now())
	if err != nil {
		return fmt.Errorf("failed to star board: %w", err)
	}
	return nil
}

// UnstarBoard removes a star from a board for a member
func (r *repository) UnstarBoard(ctx context.Context, boardID, memberID uuid.UUID) error {
	query := `DELETE FROM starred_boards WHERE id_board = $1 AND id_member = $2`
	_, err := r.conn.ExecContext(ctx, query, boardID, memberID)
	if err != nil {
		return fmt.Errorf("failed to unstar board: %w", err)
	}
	return nil
}

// IsStarred checks if a board is starred by a member
func (r *repository) IsStarred(ctx context.Context, boardID, memberID uuid.UUID) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM starred_boards WHERE id_board = $1 AND id_member = $2)`
	err := r.conn.GetContext(ctx, &exists, query, boardID, memberID)
	if err != nil {
		return false, fmt.Errorf("failed to check if board is starred: %w", err)
	}
	return exists, nil
}

// GetBoardLists retrieves all lists for a board
func (r *repository) GetBoardLists(ctx context.Context, boardID uuid.UUID) ([]*models.List, error) {
	lists := []*models.List{}
	query := `
		SELECT id, name, id_board, position, archived, created_at, updated_at
		FROM lists
		WHERE id_board = $1 AND archived = false
		ORDER BY position ASC
	`
	err := r.conn.SelectContext(ctx, &lists, query, boardID)
	if err != nil {
		return nil, fmt.Errorf("failed to get board lists: %w", err)
	}
	return lists, nil
}
