CREATE TABLE poll_votes (
    id SERIAL PRIMARY KEY,
    post_id INTEGER NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    option_id VARCHAR(255) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(post_id, user_id, option_id)
);

CREATE INDEX idx_poll_votes_post_id ON poll_votes(post_id);
CREATE INDEX idx_poll_votes_user_id ON poll_votes(user_id);
CREATE INDEX idx_poll_votes_option_id ON poll_votes(option_id);