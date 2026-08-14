-- +goose Up
CREATE TABLE posts (
  id UUID PRIMARY KEY,
  feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  published_at TIMESTAMP,
  title TEXT NOT NULL,
  url TEXT NOT NULL UNIQUE,
  description TEXT
);


-- +goose Down
DROP TABLE posts;
