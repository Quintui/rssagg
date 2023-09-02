package main

import (
	"time"

	"github.com/google/uuid"
	"github.com/quintui/rssagg/internal/database"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func databaseUserToUser(databaseUser database.User) User {
	return User{
		ID:        databaseUser.ID,
		Username:  databaseUser.Username,
		CreatedAt: databaseUser.CreatedAt,
		UpdatedAt: databaseUser.UpdatedAt,
	}
}
