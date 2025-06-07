-- name: CreateLink :one
INSERT INTO link (
    link_id, from_lume_id, to_lume_id, link_type,
    travel_details, notes, sequence_index, created_at, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8, $9
) RETURNING id, link_id, from_lume_id, to_lume_id, link_type,
    travel_details, notes, sequence_index, created_at, updated_at;

-- name: GetLinkByID :one
SELECT id, link_id, from_lume_id, to_lume_id, link_type,
    travel_details, notes, sequence_index, created_at, updated_at
FROM link WHERE id = $1;

-- name: GetLinkByLinkID :one
SELECT id, link_id, from_lume_id, to_lume_id, link_type,
    travel_details, notes, sequence_index, created_at, updated_at
FROM link WHERE link_id = $1;

-- name: ListLinksByFromLumeID :many
SELECT id, link_id, from_lume_id, to_lume_id, link_type,
    travel_details, notes, sequence_index, created_at, updated_at
FROM link 
WHERE from_lume_id = $1
ORDER BY sequence_index ASC NULLS LAST, created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListLinksByToLumeID :many
SELECT id, link_id, from_lume_id, to_lume_id, link_type,
    travel_details, notes, sequence_index, created_at, updated_at
FROM link 
WHERE to_lume_id = $1
ORDER BY sequence_index ASC NULLS LAST, created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListLinksByEitherLumeID :many
SELECT id, link_id, from_lume_id, to_lume_id, link_type,
    travel_details, notes, sequence_index, created_at, updated_at
FROM link 
WHERE from_lume_id = $1 OR to_lume_id = $1
ORDER BY sequence_index ASC NULLS LAST, created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListLinksByType :many
SELECT id, link_id, from_lume_id, to_lume_id, link_type,
    travel_details, notes, sequence_index, created_at, updated_at
FROM link 
WHERE link_type = $1
ORDER BY sequence_index ASC NULLS LAST, created_at DESC
LIMIT $2 OFFSET $3;

-- name: ListLinksByLumeIDAndType :many
SELECT id, link_id, from_lume_id, to_lume_id, link_type,
    travel_details, notes, sequence_index, created_at, updated_at
FROM link 
WHERE (from_lume_id = $1 OR to_lume_id = $1) AND link_type = $2
ORDER BY sequence_index ASC NULLS LAST, created_at DESC
LIMIT $3 OFFSET $4;

-- name: UpdateLink :one
UPDATE link SET
    from_lume_id = $2,
    to_lume_id = $3,
    link_type = $4,
    travel_details = $5,
    notes = $6,
    sequence_index = $7,
    updated_at = $8
WHERE link_id = $1
RETURNING id, link_id, from_lume_id, to_lume_id, link_type,
    travel_details, notes, sequence_index, created_at, updated_at;

-- name: DeleteLink :exec
DELETE FROM link WHERE id = $1;

-- name: DeleteLinkByLinkID :exec
DELETE FROM link WHERE link_id = $1;

-- name: CountLinksByLumeID :one
SELECT COUNT(*) FROM link WHERE from_lume_id = $1 OR to_lume_id = $1;

-- name: CountLinksByFromLumeID :one
SELECT COUNT(*) FROM link WHERE from_lume_id = $1;

-- name: CountLinksByToLumeID :one
SELECT COUNT(*) FROM link WHERE to_lume_id = $1;