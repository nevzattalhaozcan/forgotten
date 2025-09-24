BEGIN;

-- Revert to 0 (simple reversible behavior)
UPDATE users SET books_read = 0;

COMMIT;