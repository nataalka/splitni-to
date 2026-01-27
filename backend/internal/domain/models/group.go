package models

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Members   []User    `json:"members,omitempty" db:"-"`
}

type GroupBalance struct {
	UserID   uuid.UUID `json:"user_id"`
	UserName string    `json:"user_name"`
	Balance  float64   `json:"balance"`
}

type CreateGroupRequest struct {
	Name string `json:"name"`
}

type AddMemberRequest struct {
	UserID uuid.UUID `json:"user_id"`
}
