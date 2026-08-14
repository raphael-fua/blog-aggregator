-- +goose Up
CREATE TABLE feeds (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  name TEXT NOT NULL,
  url TEXT NOT NULL UNIQUE,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE
);


-- +goose Down
DROP TABLE feeds;


-- +goose Up
ALTER TABLE feeds ADD COLUMN last_fetched_at TIMESTAMP;


-- +goose Down
ALTER TABLE feeds DROP COLUMN last_fetched_at;













