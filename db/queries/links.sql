-- name: CreateLink :one
INSERT INTO links (original_url, short_name, short_url)
VALUES ($1, $2, $3)
RETURNING id, original_url, short_name, short_url, created_at, updated_at;

-- name: GetLinkByID :one
SELECT id, original_url, short_name, short_url, created_at, updated_at
FROM links
WHERE id = $1;

-- name: GetLinks :many
SELECT id, original_url, short_name, short_url, created_at, updated_at
FROM links;

-- name: UpdateLink :one
UPDATE links
SET original_url = $2,
    short_name = $3,
    short_url = $4,
    updated_at = now()
WHERE id = $1
RETURNING id, original_url, short_name, short_url, created_at, updated_at;

-- name: DeleteLink :execrows
DELETE FROM links
WHERE id = $1;
