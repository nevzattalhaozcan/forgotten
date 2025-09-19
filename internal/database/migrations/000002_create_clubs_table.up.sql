CREATE TABLE IF NOT EXISTS clubs (
  id SERIAL PRIMARY KEY,
  name VARCHAR(100) NOT NULL UNIQUE,
  description TEXT,
  location VARCHAR(255),
  genre VARCHAR(100),
  cover_image_url TEXT,
  is_private BOOLEAN NOT NULL DEFAULT FALSE,
  max_members INT NOT NULL DEFAULT 100,
  members_count INT NOT NULL DEFAULT 0,
  rating REAL NOT NULL DEFAULT 0,
  tags TEXT[],
  owner_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  current_book JSONB,
  next_meeting JSONB,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS club_memberships (
  id SERIAL PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  club_id INTEGER NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
  role VARCHAR(20) NOT NULL DEFAULT 'member',
  is_approved BOOLEAN NOT NULL DEFAULT FALSE,
  joined_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  CONSTRAINT uq_club_user UNIQUE (club_id, user_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS uni_clubs_name ON clubs (name) WHERE deleted_at IS NULL;