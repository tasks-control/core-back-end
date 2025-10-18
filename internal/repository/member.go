package repository

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"github.com/tasks-control/core-back-end/internal/models"
)

// CreateMember inserts a new member into the database
func (r *repository) CreateMember(ctx context.Context, member *models.Member) error {
	query := `
		INSERT INTO members (id, email, username, full_name, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.conn.ExecContext(ctx, query,
		member.ID,
		member.Email,
		member.Username,
		member.FullName,
		member.PasswordHash,
		member.CreatedAt,
		member.UpdatedAt,
	)
	return err
}

// GetMemberByEmail retrieves a member by email
func (r *repository) GetMemberByEmail(ctx context.Context, email string) (*models.Member, error) {
	var member models.Member
	query := `
		SELECT id, email, username, full_name, password_hash, created_at, updated_at
		FROM members
		WHERE email = $1
	`
	err := r.conn.GetContext(ctx, &member, query, email)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// GetMemberByID retrieves a member by ID
func (r *repository) GetMemberByID(ctx context.Context, id uuid.UUID) (*models.Member, error) {
	var member models.Member
	query := `
		SELECT id, email, username, full_name, password_hash, created_at, updated_at
		FROM members
		WHERE id = $1
	`
	err := r.conn.GetContext(ctx, &member, query, id)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// GetMemberByUsername retrieves a member by username
func (r *repository) GetMemberByUsername(ctx context.Context, username string) (*models.Member, error) {
	var member models.Member
	query := `
		SELECT id, email, username, full_name, password_hash, created_at, updated_at
		FROM members
		WHERE username = $1
	`
	err := r.conn.GetContext(ctx, &member, query, username)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &member, nil
}

// UpdateMember updates an existing member
func (r *repository) UpdateMember(ctx context.Context, member *models.Member) error {
	query := `
		UPDATE members
		SET email = $2, username = $3, full_name = $4, password_hash = $5, updated_at = $6
		WHERE id = $1
	`
	_, err := r.conn.ExecContext(ctx, query,
		member.ID,
		member.Email,
		member.Username,
		member.FullName,
		member.PasswordHash,
		member.UpdatedAt,
	)
	return err
}
