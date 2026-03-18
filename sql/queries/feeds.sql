-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id, last_fetched_at)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
ON CONFLICT (url) DO UPDATE
SET url = EXCLUDED.url
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feeds;

-- name: GetAllFeedsWithUser :many
SELECT feeds.name, feeds.url, users.name AS user_name
FROM feeds
INNER JOIN users ON feeds.user_id = users.id;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1;