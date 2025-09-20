-- name: GetUser :one
SELECT name FROM users
WHERE id = $1 LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (name, password)
VALUES ($1, $2)
RETURNING *;
