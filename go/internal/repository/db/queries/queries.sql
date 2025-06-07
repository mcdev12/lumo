-- name: CreateLume :one
INSERT INTO lume (lume_id, lumo_id, label, type, description, metadata, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)  -- $1 = lume_id (provided by Go!)
RETURNING id, lume_id, lumo_id, label, type, description, metadata, created_at, updated_at;

-- name: GetLumeByID :one  
SELECT id, lume_id, lumo_id, label, type, description, metadata, created_at, updated_at
FROM lume WHERE id = $1;

-- name: GetLumeByLumeID :one
SELECT id, lume_id, lumo_id, label, type, description, metadata, created_at, updated_at
FROM lume WHERE lume_id = $1;

-- name: ListLumesByLumoID :many
SELECT id, lume_id, lumo_id, label, type, description, metadata, created_at, updated_at
FROM lume 
WHERE lumo_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: UpdateLume :one
UPDATE lume 
SET label = $2, type = $3, description = $4, metadata = $5, updated_at = $6
WHERE id = $1
RETURNING id, lume_id, lumo_id, label, type, description, metadata, created_at, updated_at;

-- name: DeleteLume :exec
DELETE FROM lume WHERE id = $1;