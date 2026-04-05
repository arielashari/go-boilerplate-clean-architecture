ALTER TABLE users DROP CONSTRAINT IF EXISTS fk_users_role;

ALTER TABLE users DROP COLUMN IF EXISTS role_id;

DELETE FROM roles WHERE name IN ('admin', 'user');