// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: feed_posts.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createPost = `-- name: CreatePost :one
INSERT INTO posts(id, title, url, description, feed_id, created_at, updated_at, published_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id, title, url, description, published_at, created_at, updated_at, feed_id
`

type CreatePostParams struct {
	ID          uuid.UUID
	Title       string
	Url         string
	Description sql.NullString
	FeedID      uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	PublishedAt sql.NullTime
}

func (q *Queries) CreatePost(ctx context.Context, arg CreatePostParams) (Post, error) {
	row := q.db.QueryRowContext(ctx, createPost,
		arg.ID,
		arg.Title,
		arg.Url,
		arg.Description,
		arg.FeedID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.PublishedAt,
	)
	var i Post
	err := row.Scan(
		&i.ID,
		&i.Title,
		&i.Url,
		&i.Description,
		&i.PublishedAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FeedID,
	)
	return i, err
}

const getPostsForUser = `-- name: GetPostsForUser :many
SELECT posts.id, posts.title, posts.url, posts.description, posts.published_at, posts.created_at, posts.updated_at, posts.feed_id from posts 
JOIN feed_follows ON feed_follows.feed_id = posts.feed_id
WHERE feed_follows.user_id = $1
ORDER BY posts.published_at DESC
LIMIT $2
`

type GetPostsForUserParams struct {
	UserID uuid.UUID
	Limit  int32
}

func (q *Queries) GetPostsForUser(ctx context.Context, arg GetPostsForUserParams) ([]Post, error) {
	rows, err := q.db.QueryContext(ctx, getPostsForUser, arg.UserID, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Post
	for rows.Next() {
		var i Post
		if err := rows.Scan(
			&i.ID,
			&i.Title,
			&i.Url,
			&i.Description,
			&i.PublishedAt,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FeedID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}