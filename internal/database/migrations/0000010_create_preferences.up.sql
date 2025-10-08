ALTER TABLE users ADD COLUMN preferences JSONB DEFAULT '{}';

-- Create index for faster preferences queries
CREATE INDEX idx_users_preferences ON users USING GIN (preferences);