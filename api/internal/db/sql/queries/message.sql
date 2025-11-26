-- name: CreateMessage :exec
INSERT INTO messages (room_id, author_id, content)
VALUES ($1, $2, $3);

-- name: GetMessagesByRoomID :many
SELECT
    m.id AS message_id,
    m.room_id,
    m.content,
    m.created_at,
    m.updated_at,
    m.author_id,
    a.username AS author_name
FROM messages AS m
INNER JOIN accounts AS a ON m.author_id = a.id
WHERE m.room_id = $1
ORDER BY m.created_at ASC;
