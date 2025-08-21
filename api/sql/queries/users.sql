-- name: GetUserByID :one
SELECT id,
    email,
    name,
    created_at,
    updated_at
FROM users
WHERE id = $1
    AND deleted_at IS NULL;

-- name: ListUsers :many
SELECT id,
    email,
    name,
    created_at,
    updated_at
FROM users
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;