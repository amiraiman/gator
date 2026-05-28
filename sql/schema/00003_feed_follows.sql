-- +goose Up
CREATE TABLE feed_follows(
    id UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
    CONSTRAINT unique_user_and_feed UNIQUE (user_id, feed_id)
);

INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
SELECT gen_random_uuid(), NOW(), NOW(), user_id, id
FROM feeds;

-- +goose Down
DROP TABLE feed_follows;
