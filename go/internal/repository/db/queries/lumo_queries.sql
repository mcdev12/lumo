-- name: CreateLumo :one
INSERT INTO lumo (
    lumo_id, user_id, title, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING id, lumo_id, user_id, title, created_at, updated_at;

-- name: GetLumoByID :one
SELECT id, lumo_id, user_id, title, created_at, updated_at
FROM lumo WHERE id = $1;

-- name: GetLumoByLumoID :one
SELECT id, lumo_id, user_id, title, created_at, updated_at
FROM lumo WHERE lumo_id = $1;

-- name: ListLumosByUserID :many
SELECT id, lumo_id, user_id, title, created_at, updated_at
FROM lumo 
WHERE user_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateLumo :one
UPDATE lumo SET
    title = $2,
    updated_at = $3
WHERE lumo_id = $1
RETURNING id, lumo_id, user_id, title, created_at, updated_at;

-- name: DeleteLumo :exec
DELETE FROM lumo WHERE id = $1;

-- name: DeleteLumoByLumoID :exec
DELETE FROM lumo WHERE lumo_id = $1;

-- name: CountLumosByUserID :one
SELECT COUNT(*) FROM lumo WHERE user_id = $1;