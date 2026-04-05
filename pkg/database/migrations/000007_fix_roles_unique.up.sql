ALTER TABLE roles ADD CONSTRAINT roles_name_unique UNIQUE (name);

INSERT INTO roles (id, name) VALUES (gen_random_uuid(), 'admin') ON CONFLICT (name) DO NOTHING;