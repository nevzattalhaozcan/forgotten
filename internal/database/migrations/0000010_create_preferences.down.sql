ALTER TABLE users DROP COLUMN IF EXISTS preferences;
DROP INDEX IF EXISTS idx_users_preferences;