-- name: CreateAccount :one
INSERT INTO accounts (username, password_hash)
VALUES ($1, $2)
RETURNING id;

-- name: GetAccountByID :one
SELECT id, username, created_at, updated_at
FROM accounts
WHERE id = $1;

-- name: GetAccountByUsername :one
SELECT id, username, created_at, updated_at
FROM accounts
WHERE username = $1;

-- name: GetLoginCredential :one
SELECT id, username, password_hash
FROM accounts
WHERE username = $1;
