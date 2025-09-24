BEGIN;
DROP TRIGGER IF EXISTS trg_club_assignments_updated ON club_book_assignments;
DROP TRIGGER IF EXISTS trg_user_book_progress_updated ON user_book_progress;
DROP FUNCTION IF EXISTS set_updated_at();

DROP TABLE IF EXISTS club_book_assignments;
DROP TABLE IF EXISTS reading_logs;
DROP TABLE IF EXISTS user_book_progress;
COMMIT;