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

-- name: GetFeedFromID :one
SELECT * FROM feeds
  WHERE id = $1;

-- name: MarkFeedFetch :exec
UPDATE feeds 
  SET updated_at = $1, last_fetched_at=$1
WHERE id = $2;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
  ORDER BY last_fetched_at ASC NULLS FIRST
  LIMIT 1;
