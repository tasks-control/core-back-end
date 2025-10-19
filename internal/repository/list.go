package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tasks-control/core-back-end/internal/models"
)

// CreateList inserts a new list into the database
func (r *repository) CreateList(ctx context.Context, list *models.List) error {
	query := `
		INSERT INTO lists (id, name, id_board, position, archived, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.conn.ExecContext(ctx, query,
		list.ID,
		list.Name,
		list.IDBoard,
		list.Position,
		list.Archived,
		list.CreatedAt,
		list.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create list: %w", err)
	}
	return nil
}

// GetListByID retrieves a list by ID
func (r *repository) GetListByID(ctx context.Context, listID uuid.UUID) (*models.List, error) {
	var list models.List
	query := `
		SELECT id, name, id_board, position, archived, created_at, updated_at
		FROM lists
		WHERE id = $1
	`
	err := r.conn.GetContext(ctx, &list, query, listID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get list: %w", err)
	}
	return &list, nil
}

// UpdateList updates an existing list
func (r *repository) UpdateList(ctx context.Context, list *models.List) error {
	query := `
		UPDATE lists
		SET name = $2, position = $3, archived = $4, updated_at = $5
		WHERE id = $1
	`
	result, err := r.conn.ExecContext(ctx, query,
		list.ID,
		list.Name,
		list.Position,
		list.Archived,
		list.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update list: %w", err)
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

// DeleteList deletes a list (cascade will delete related cards)
func (r *repository) DeleteList(ctx context.Context, listID uuid.UUID) error {
	query := `DELETE FROM lists WHERE id = $1`
	result, err := r.conn.ExecContext(ctx, query, listID)
	if err != nil {
		return fmt.Errorf("failed to delete list: %w", err)
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

// GetListCards retrieves all cards in a list
func (r *repository) GetListCards(ctx context.Context, listID uuid.UUID) ([]*models.Card, error) {
	cards := []*models.Card{}
	query := `
		SELECT id, title, description, id_list, position, archived, created_by, created_at, updated_at
		FROM cards
		WHERE id_list = $1 AND archived = false
		ORDER BY position ASC
	`
	err := r.conn.SelectContext(ctx, &cards, query, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to get list cards: %w", err)
	}
	return cards, nil
}

// GetMaxListPosition returns the maximum position value for lists in a board
func (r *repository) GetMaxListPosition(ctx context.Context, boardID uuid.UUID) (float64, error) {
	var maxPosition sql.NullFloat64
	query := `
		SELECT MAX(position)
		FROM lists
		WHERE id_board = $1 AND archived = false
	`
	err := r.conn.GetContext(ctx, &maxPosition, query, boardID)
	if err != nil {
		return 0, fmt.Errorf("failed to get max list position: %w", err)
	}

	if !maxPosition.Valid {
		return 0, nil
	}

	return maxPosition.Float64, nil
}

// GetListCountInBoard returns the number of lists in a board
func (r *repository) GetListCountInBoard(ctx context.Context, boardID uuid.UUID) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM lists
		WHERE id_board = $1 AND archived = false
	`
	err := r.conn.GetContext(ctx, &count, query, boardID)
	if err != nil {
		return 0, fmt.Errorf("failed to get list count: %w", err)
	}
	return count, nil
}
