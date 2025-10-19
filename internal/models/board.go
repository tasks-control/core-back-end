package models

import (
	"time"

	"github.com/google/uuid"
)

// Board represents a task board in the system
type Board struct {
	ID              uuid.UUID `db:"id" json:"id"`
	Name            string    `db:"name" json:"name"`
	NameBoardUnique string    `db:"name_board_unique" json:"name_board_unique"`
	Description     *string   `db:"description" json:"description,omitempty"`
	PasswordHash    string    `db:"password_hash" json:"-"`
	IDMemberCreator uuid.UUID `db:"id_member_creator" json:"idMemberCreator"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt       time.Time `db:"updated_at" json:"updatedAt"`
	Starred         *bool     `db:"starred" json:"starred,omitempty"`          // Only populated in list queries
	MemberCount     *int      `db:"member_count" json:"memberCount,omitempty"` // Only populated in list queries
}

// BoardMember represents the relationship between a board and a member
type BoardMember struct {
	ID       uuid.UUID `db:"id" json:"id"`
	IDBoard  uuid.UUID `db:"id_board" json:"idBoard"`
	IDMember uuid.UUID `db:"id_member" json:"idMember"`
	Role     string    `db:"role" json:"role"`
	JoinedAt time.Time `db:"joined_at" json:"joinedAt"`
}

// StarredBoard represents a starred board relationship
type StarredBoard struct {
	ID        uuid.UUID `db:"id" json:"id"`
	IDBoard   uuid.UUID `db:"id_board" json:"idBoard"`
	IDMember  uuid.UUID `db:"id_member" json:"idMember"`
	StarredAt time.Time `db:"starred_at" json:"starredAt"`
}

// List represents a list within a board
type List struct {
	ID        uuid.UUID `db:"id" json:"id"`
	Name      string    `db:"name" json:"name"`
	IDBoard   uuid.UUID `db:"id_board" json:"idBoard"`
	Position  float64   `db:"position" json:"position"`
	Archived  bool      `db:"archived" json:"archived"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}

// Card represents a task card within a list
type Card struct {
	ID          uuid.UUID `db:"id" json:"id"`
	Title       string    `db:"title" json:"title"`
	Description *string   `db:"description" json:"description,omitempty"`
	IDList      uuid.UUID `db:"id_list" json:"idList"`
	Position    float64   `db:"position" json:"position"`
	Archived    bool      `db:"archived" json:"archived"`
	CreatedBy   uuid.UUID `db:"created_by" json:"createdBy"`
	CreatedAt   time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time `db:"updated_at" json:"updatedAt"`
}

// BoardRole constants
const (
	BoardRoleOwner     = "owner"
	BoardRoleModerator = "moderator"
	BoardRoleMember    = "member"
)
