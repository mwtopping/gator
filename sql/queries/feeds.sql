-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
)
RETURNING *;

-- name: ListFeeds :many
SELECT * FROM feeds;

-- name: GetFeedUser :one
SELECT name FROM users
  WHERE id = $1;

-- name: GetFeedFromURL :one
SELECT * FROM feeds
  WHERE url = $1;
