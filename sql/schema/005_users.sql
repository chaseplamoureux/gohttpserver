-- +goose Up
ALTER TABLE users
ADD is_chirpy_red BOOLEAN;

-- +goose Down
ALTER TABLE users ALTER COLUMN is_chirpy_red 
SET DEFAULT FALSE;