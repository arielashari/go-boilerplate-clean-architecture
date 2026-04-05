INSERT INTO roles (id, name)
VALUES 
    (gen_random_uuid(), 'admin'),
    (gen_random_uuid(), 'user')
ON CONFLICT (name) DO NOTHING;

ALTER TABLE users ADD COLUMN role_id UUID;

UPDATE users 
SET role_id = (SELECT id FROM roles WHERE name = 'user' LIMIT 1)
WHERE role_id IS NULL;

ALTER TABLE users ALTER COLUMN role_id SET NOT NULL;

ALTER TABLE users 
ADD CONSTRAINT fk_users_role 
FOREIGN KEY (role_id) REFERENCES roles(id);