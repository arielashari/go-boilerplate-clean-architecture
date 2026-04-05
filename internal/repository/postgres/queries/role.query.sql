-- name: CreateRole :one
INSERT INTO roles (
    id, name
) VALUES (
    $1, $2
) RETURNING *;

-- name: GetRoleByID :one
SELECT * FROM roles
WHERE id = $1 AND deleted_at IS NULL LIMIT 1;

-- name: ListRoles :many
SELECT * FROM roles
WHERE deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: CountRoles :one
SELECT COUNT(*) FROM roles
WHERE deleted_at IS NULL;   

-- name: DeleteRole :execresult
UPDATE roles
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;   

-- name: UpdateRole :one
UPDATE roles
SET name = $2, 
    updated_at = NOW()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;