package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id" postgres:"id"`
	Name         string    `json:"name" postgres:"name"`
	Surname      string    `json:"surname" postgres:"surname"`
	Email        string    `json:"email" postgres:"email"`
	PasswordHash string    `json:"-" postgres:"password_hash"`
	CreatedAt    time.Time `json:"created_at" postgres:"created_at"`
}

type CreateUser struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
