-- +goose Up
ALTER TABLE posts ADD CONSTRAINT posts_url_unique UNIQUE (url);

-- +goose Down
ALTER TABLE posts DROP CONSTRAINT posts_url_unique;
