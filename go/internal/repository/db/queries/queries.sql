-- name: CreateLume :one
INSERT INTO lume (
    lume_id, lumo_id, type, name,
    date_start, date_end, latitude, longitude,
    address, description, images, category_tags,
    booking_link, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15
) RETURNING id, lume_id, lumo_id, type, name,
    date_start, date_end, latitude, longitude,
    address, description, images, category_tags,
    booking_link, created_at, updated_at;

-- name: GetLumeByID :one
SELECT id, lume_id, lumo_id, type, name,
    date_start, date_end, latitude, longitude,
    address, description, images, category_tags,
    booking_link, created_at, updated_at
FROM lume WHERE id = $1;

-- name: GetLumeByLumeID :one
SELECT id, lume_id, lumo_id, type, name,
    date_start, date_end, latitude, longitude,
    address, description, images, category_tags,
    booking_link, created_at, updated_at
FROM lume WHERE lume_id = $1;

-- name: ListLumesByLumoID :many
SELECT id, lume_id, lumo_id, type, name,
    date_start, date_end, latitude, longitude,
    address, description, images, category_tags,
    booking_link, created_at, updated_at
FROM lume 
WHERE lumo_id = $1
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListLumesByType :many
SELECT id, lume_id, lumo_id, type, name,
    date_start, date_end, latitude, longitude,
    address, description, images, category_tags,
    booking_link, created_at, updated_at
FROM lume 
WHERE lumo_id = $1 AND type = $2
ORDER BY created_at DESC
LIMIT $3 OFFSET $4;

-- name: SearchLumesByLocation :many
SELECT id, lume_id, lumo_id, type, name,
    date_start, date_end, latitude, longitude,
    address, description, images, category_tags,
    booking_link, created_at, updated_at
FROM lume 
WHERE lumo_id = $1 
    AND latitude IS NOT NULL 
    AND longitude IS NOT NULL
    AND latitude BETWEEN $2 AND $3
    AND longitude BETWEEN $4 AND $5
ORDER BY created_at DESC
LIMIT $6 OFFSET $7;

-- name: UpdateLume :one
UPDATE lume SET
    name = $2,
    type = $3,
    date_start = $4,
    date_end = $5,
    latitude = $6,
    longitude = $7,
    address = $8,
    description = $9,
    images = $10,
    category_tags = $11,
    booking_link = $12,
    updated_at = $13
WHERE lume_id = $1
RETURNING id, lume_id, lumo_id, type, name,
    date_start, date_end, latitude, longitude,
    address, description, images, category_tags,
    booking_link, created_at, updated_at;

-- name: DeleteLume :exec
DELETE FROM lume WHERE id = $1;

-- name: DeleteLumeByLumeID :exec
DELETE FROM lume WHERE lume_id = $1;

-- name: CountLumesByLumo :one
SELECT COUNT(*) FROM lume WHERE lumo_id = $1;