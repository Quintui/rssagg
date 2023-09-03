// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package database

import (
	"time"

	"github.com/google/uuid"
)

type Feed struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Url       string
	AuthorID  uuid.UUID
}

type User struct {
	ID        uuid.UUID
	Username  string
	CreatedAt time.Time
	UpdatedAt time.Time
	ApiKey    string
}
