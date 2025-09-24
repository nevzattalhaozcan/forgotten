BEGIN;

-- Track per-user progress on a book
CREATE TABLE IF NOT EXISTS user_book_progress (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'reading', -- not_started|reading|paused|finished
    current_page INT,
    percent NUMERIC(5,2), -- denorm
    started_at TIMESTAMPTZ DEFAULT NOW(),
    finished_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ DEFAULT NOW(),
    UNIQUE (user_id, book_id)
);

CREATE INDEX IF NOT EXISTS idx_user_book_progress_user ON user_book_progress(user_id);
CREATE INDEX IF NOT EXISTS idx_user_book_progress_book ON user_book_progress(book_id);

-- stores each update with deltas for history and reporting
CREATE TABLE IF NOT EXISTS reading_logs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    book_id BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    club_id BIGINT REFERENCES clubs(id) ON DELETE SET NULL,
    assignment_id BIGINT, -- forward ref; not FK to avoid circular on create order
    pages_delta INT,      -- how many pages were read in this session
    from_page INT,
    to_page INT,
    minutes INT,
    note TEXT,
    created_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_reading_logs_user ON reading_logs(user_id);
CREATE INDEX IF NOT EXISTS idx_reading_logs_book ON reading_logs(book_id);
CREATE INDEX IF NOT EXISTS idx_reading_logs_club ON reading_logs(club_id);

-- Club-level assignments (what the club is reading now / history)
CREATE TABLE IF NOT EXISTS club_book_assignments (
    id BIGSERIAL PRIMARY KEY,
    club_id BIGINT NOT NULL REFERENCES clubs(id) ON DELETE CASCADE,
    book_id BIGINT NOT NULL REFERENCES books(id) ON DELETE CASCADE,
    status TEXT NOT NULL DEFAULT 'active', -- active|completed|archived
    start_date DATE DEFAULT CURRENT_DATE,
    due_date DATE,
    completed_at TIMESTAMPTZ,
    target_page INT,       -- "checkpoint" page for the club
    checkpoint TEXT,       -- description
    created_at TIMESTAMPTZ DEFAULT NOW(),
    updated_at TIMESTAMPTZ DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_club_assignments_club ON club_book_assignments(club_id);
CREATE INDEX IF NOT EXISTS idx_club_assignments_book ON club_book_assignments(book_id);

-- trigger to keep updated_at current
CREATE OR REPLACE FUNCTION set_updated_at() RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_user_book_progress_updated ON user_book_progress;
CREATE TRIGGER trg_user_book_progress_updated BEFORE UPDATE ON user_book_progress
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

DROP TRIGGER IF EXISTS trg_club_assignments_updated ON club_book_assignments;
CREATE TRIGGER trg_club_assignments_updated BEFORE UPDATE ON club_book_assignments
FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMIT;