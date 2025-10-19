package models

import (
	"time"

	"github.com/google/uuid"
)

// Member represents a user in the system
type Member struct {
	ID           uuid.UUID `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	Username     string    `db:"username" json:"username"`
	FullName     *string   `db:"full_name" json:"fullName,omitempty"`
	PasswordHash string    `db:"password_hash" json:"-"`
	CreatedAt    time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time `db:"updated_at" json:"updatedAt"`
}

// RefreshToken represents a JWT refresh token stored in the database
type RefreshToken struct {
	ID        uuid.UUID `db:"id" json:"id"`
	IDMember  uuid.UUID `db:"id_member" json:"idMember"`
	TokenHash string    `db:"token_hash" json:"-"`
	ExpiresAt time.Time `db:"expires_at" json:"expiresAt"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	Revoked   bool      `db:"revoked" json:"revoked"`
}
