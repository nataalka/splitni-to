package models

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Members   []User    `json:"members,omitempty"`
}

type GroupBalance struct {
	UserID   uuid.UUID `json:"user_id"`
	UserName string    `json:"user_name"`
	Balance  float64   `json:"balance"`
}
