-- name: CreateFeed :one
INSERT INTO feed (id, name, author_id, url, created_at, updated_at )
VALUES ($1, $2, $3, $4, $5, $6 ) 
RETURNING *;

