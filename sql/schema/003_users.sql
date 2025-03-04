-- +goose Up
ALTER TABLE users
ADD hashed_password TEXT;

ALTER TABLE users ALTER COLUMN hashed_password 
SET DEFAULT 'unset';