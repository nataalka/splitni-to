package models

import "github.com/google/uuid"

type Friendship struct {
	UserID1 uuid.UUID `json:"user_id1" db:"user_id1"`
	UserID2 uuid.UUID `json:"user_id2" db:"user_id2"`
	Status  string    `json:"status"   db:"status"`
}

type FriendRequest struct {
	Email string `json:"email"`
}
