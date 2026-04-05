ALTER TABLE users ADD COLUMN IF NOT EXISTS password VARCHAR(255);
UPDATE users SET password = 'temporary_placeholder' WHERE password IS NULL;
ALTER TABLE users ALTER COLUMN password SET NOT NULL;