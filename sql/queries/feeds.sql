-- name: GetFeedByName :one
SELECT * FROM feeds
WHERE name = $1 LIMIT 1;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1 LIMIT 1;

-- name: GetFeeds :many
SELECT * FROM feeds
ORDER BY name;

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

-- name: DeleteFeeds :exec
DELETE FROM feeds;

-- name: MarkFeedFetchedById :exec
UPDATE feeds SET last_fetched_at = $1, updated_at = $2 WHERE id = $3;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at NULLS FIRST LIMIT 1;

-- name: CreateFeedFollow :one
WITH inserted_follow_fields AS (
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
    inserted_follow_fields.*,
    users.name as user_name,
    feeds.name as feed_name
FROM inserted_follow_fields
INNER JOIN users ON users.id = inserted_follow_fields.user_id
INNER JOIN feeds ON feeds.id = inserted_follow_fields.feed_id;

-- name: DeleteFollowFeeds :exec
DELETE FROM feed_follows;

-- name: GetFeedFollowsForUser :many
SELECT
    feed_follows.*,
    users.name as user_name,
    feeds.name as feed_name
FROM feed_follows
INNER JOIN users ON users.id = feed_follows.user_id
INNER JOIN feeds ON feeds.id = feed_follows.feed_id
WHERE feed_follows.user_id = $1;

-- name: DeleteFollowFeedsByUserAndFeedId :exec
DELETE FROM feed_follows WHERE user_id = $1 AND feed_id = $2;
