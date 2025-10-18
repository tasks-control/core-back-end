package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/google/uuid"
	"github.com/tasks-control/core-back-end/internal/models"
)

// CreateRefreshToken inserts a new refresh token into the database
func (r *repository) CreateRefreshToken(ctx context.Context, token *models.RefreshToken) error {
	query := `
		INSERT INTO refresh_tokens (id, id_member, token_hash, expires_at, created_at, revoked)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.conn.ExecContext(ctx, query,
		token.ID,
		token.IDMember,
		token.TokenHash,
		token.ExpiresAt,
		token.CreatedAt,
		token.Revoked,
	)
	return err
}

// GetRefreshTokenByHash retrieves a refresh token by its hash
func (r *repository) GetRefreshTokenByHash(ctx context.Context, tokenHash string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	query := `
		SELECT id, id_member, token_hash, expires_at, created_at, revoked
		FROM refresh_tokens
		WHERE token_hash = $1 AND revoked = false AND expires_at > NOW()
	`
	err := r.conn.GetContext(ctx, &token, query, tokenHash)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &token, nil
}

// RevokeRefreshToken marks a refresh token as revoked
func (r *repository) RevokeRefreshToken(ctx context.Context, tokenHash string) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = true
		WHERE token_hash = $1
	`
	_, err := r.conn.ExecContext(ctx, query, tokenHash)
	return err
}

// RevokeAllUserTokens revokes all refresh tokens for a user
func (r *repository) RevokeAllUserTokens(ctx context.Context, userID uuid.UUID) error {
	query := `
		UPDATE refresh_tokens
		SET revoked = true
		WHERE id_member = $1 AND revoked = false
	`
	_, err := r.conn.ExecContext(ctx, query, userID)
	return err
}

// DeleteExpiredTokens removes expired tokens from the database
func (r *repository) DeleteExpiredTokens(ctx context.Context) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE expires_at < NOW() OR (revoked = true AND created_at < NOW() - INTERVAL '30 days')
	`
	_, err := r.conn.ExecContext(ctx, query)
	return err
}
