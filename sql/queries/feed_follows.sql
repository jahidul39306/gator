-- name: CreateFeedFollow :one
WITH inserted AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES (
        $1,
        $2,
        $3,
        $4,
        $5
    )
    RETURNING *
)
SELECT
    inserted.*,
    users.name AS user_name,
    feeds.name As feed_name
FROM inserted
INNER JOIN users ON inserted.user_id = users.id
INNER JOIN feeds ON inserted.feed_id = feeds.id;