-- name: CreateRoom :exec
INSERT INTO rooms (name, created_by)
VALUES ($1, $2);

-- name: GetRooms :many
SELECT id, name, created_by, created_at, updated_at
FROM rooms
ORDER BY updated_at;
