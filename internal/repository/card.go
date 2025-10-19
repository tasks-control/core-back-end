package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/tasks-control/core-back-end/internal/models"
)

// CreateCard inserts a new card into the database
func (r *repository) CreateCard(ctx context.Context, card *models.Card) error {
	query := `
		INSERT INTO cards (id, title, description, id_list, position, archived, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.conn.ExecContext(ctx, query,
		card.ID,
		card.Title,
		card.Description,
		card.IDList,
		card.Position,
		card.Archived,
		card.CreatedBy,
		card.CreatedAt,
		card.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create card: %w", err)
	}
	return nil
}

// GetCardByID retrieves a card by ID
func (r *repository) GetCardByID(ctx context.Context, cardID uuid.UUID) (*models.Card, error) {
	var card models.Card
	query := `
		SELECT id, title, description, id_list, position, archived, created_by, created_at, updated_at
		FROM cards
		WHERE id = $1
	`
	err := r.conn.GetContext(ctx, &card, query, cardID)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get card: %w", err)
	}
	return &card, nil
}

// UpdateCard updates an existing card
func (r *repository) UpdateCard(ctx context.Context, card *models.Card) error {
	query := `
		UPDATE cards
		SET title = $2, description = $3, id_list = $4, position = $5, archived = $6, updated_at = $7
		WHERE id = $1
	`
	result, err := r.conn.ExecContext(ctx, query,
		card.ID,
		card.Title,
		card.Description,
		card.IDList,
		card.Position,
		card.Archived,
		card.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to update card: %w", err)
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

// DeleteCard deletes a card
func (r *repository) DeleteCard(ctx context.Context, cardID uuid.UUID) error {
	query := `DELETE FROM cards WHERE id = $1`
	result, err := r.conn.ExecContext(ctx, query, cardID)
	if err != nil {
		return fmt.Errorf("failed to delete card: %w", err)
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

// GetMaxCardPosition returns the maximum position value for cards in a list
func (r *repository) GetMaxCardPosition(ctx context.Context, listID uuid.UUID) (float64, error) {
	var maxPosition sql.NullFloat64
	query := `
		SELECT MAX(position)
		FROM cards
		WHERE id_list = $1 AND archived = false
	`
	err := r.conn.GetContext(ctx, &maxPosition, query, listID)
	if err != nil {
		return 0, fmt.Errorf("failed to get max card position: %w", err)
	}

	if !maxPosition.Valid {
		return 0, nil
	}

	return maxPosition.Float64, nil
}

// GetCardCountInList returns the number of cards in a list
func (r *repository) GetCardCountInList(ctx context.Context, listID uuid.UUID) (int, error) {
	var count int
	query := `
		SELECT COUNT(*)
		FROM cards
		WHERE id_list = $1 AND archived = false
	`
	err := r.conn.GetContext(ctx, &count, query, listID)
	if err != nil {
		return 0, fmt.Errorf("failed to get card count: %w", err)
	}
	return count, nil
}

// GetBoardIDByCardID retrieves the board ID for a given card (via list)
func (r *repository) GetBoardIDByCardID(ctx context.Context, cardID uuid.UUID) (uuid.UUID, error) {
	var boardID uuid.UUID
	query := `
		SELECT l.id_board
		FROM cards c
		INNER JOIN lists l ON c.id_list = l.id
		WHERE c.id = $1
	`
	err := r.conn.GetContext(ctx, &boardID, query, cardID)
	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, nil
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get board ID by card: %w", err)
	}
	return boardID, nil
}

// GetBoardIDByListID retrieves the board ID for a given list
func (r *repository) GetBoardIDByListID(ctx context.Context, listID uuid.UUID) (uuid.UUID, error) {
	var boardID uuid.UUID
	query := `
		SELECT id_board
		FROM lists
		WHERE id = $1
	`
	err := r.conn.GetContext(ctx, &boardID, query, listID)
	if errors.Is(err, sql.ErrNoRows) {
		return uuid.Nil, nil
	}
	if err != nil {
		return uuid.Nil, fmt.Errorf("failed to get board ID by list: %w", err)
	}
	return boardID, nil
}
