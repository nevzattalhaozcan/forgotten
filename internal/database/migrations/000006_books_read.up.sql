BEGIN;

-- Set books_read based on finished entries in user_book_progress
UPDATE users u
SET books_read = COALESCE(t.cnt, 0)
FROM (
  SELECT user_id, COUNT(*) AS cnt
  FROM user_book_progress
  WHERE status = 'finished'
  GROUP BY user_id
) t
WHERE u.id = t.user_id;

-- Ensure users with no finished entries have books_read = 0
UPDATE users SET books_read = 0 WHERE books_read IS NULL;

COMMIT;