package models

import (
	"time"

	"github.com/google/uuid"
)

type Group struct {
	ID        uuid.UUID `json:"id" postgres:"id"`
	Name      string    `json:"name" postgres:"name"`
	CreatedAt time.Time `json:"created_at" postgres:"created_at"`
	Members   []User    `json:"members,omitempty"`
}

type GroupBalance struct {
	UserID   uuid.UUID `json:"user_id"`
	UserName string    `json:"user_name"`
	Balance  float64   `json:"balance"`
}
