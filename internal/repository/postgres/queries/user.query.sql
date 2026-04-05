-- name: CreateUser :one
INSERT INTO users (
    id, email, password, first_name, last_name, phone_prefix, phone_number, role_id
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, $8
) RETURNING *;

-- name: GetUserByID :one
SELECT sqlc.embed(users), 
    sqlc.embed(roles) FROM users
LEFT JOIN roles ON users.role_id = roles.id
WHERE users.id = $1 AND users.deleted_at IS NULL LIMIT 1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1 AND deleted_at IS NULL LIMIT 1;

-- name: GetUserForAuth :one
SELECT id,email, password FROM users
WHERE email = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListUsers :many
SELECT 
    sqlc.embed(users), 
    sqlc.embed(roles)
FROM users
LEFT JOIN roles ON users.role_id = roles.id
WHERE 
    (users.deleted_at IS NULL)
    AND (users.role_id = sqlc.narg('role_id')::uuid OR sqlc.narg('role_id') IS NULL)
    AND (
        sqlc.narg('search')::text IS NULL OR 
        users.first_name ILIKE '%' || sqlc.narg('search') || '%' OR 
        users.last_name ILIKE '%' || sqlc.narg('search') || '%' OR 
        users.email ILIKE '%' || sqlc.narg('search') || '%'
    )
ORDER BY 
    CASE WHEN sqlc.arg('sort_by') = 'first_name' AND sqlc.arg('sort_dir') = 'asc' THEN users.first_name END ASC,
    CASE WHEN sqlc.arg('sort_by') = 'first_name' AND sqlc.arg('sort_dir') = 'desc' THEN users.first_name END DESC,
    CASE WHEN sqlc.arg('sort_by') = 'created_at' AND sqlc.arg('sort_dir') = 'asc' THEN users.created_at END ASC,
    users.created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountUsers :one
SELECT COUNT(*) FROM users
WHERE deleted_at IS NULL;

-- name: DeleteUser :execresult
UPDATE users
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: UpdateUser :one
UPDATE users
SET email = $2, 
    first_name = $3, 
    last_name = $4, 
    phone_prefix = $5, 
    phone_number = $6, 
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;