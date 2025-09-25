-- name: GetUser :one
SELECT id, name, password FROM users
WHERE id = $1 LIMIT 1;

-- name: GetUserByName :one
SELECT id, name, password FROM users
WHERE name = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (name, password)
VALUES ($1, $2)
RETURNING *;

-- name: CreateUrl :one
INSERT INTO urls (user_id, orginal_url, short_url)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetUserURLs :many
SELECT id, user_id, orginal_url, short_url FROM urls
WHERE user_id = $1;
