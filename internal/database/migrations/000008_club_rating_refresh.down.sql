BEGIN;

-- Drop triggers on the club_ratings table
DROP TRIGGER IF EXISTS trg_club_ratings_ins ON club_ratings;
DROP TRIGGER IF EXISTS trg_club_ratings_upd ON club_ratings;
DROP TRIGGER IF EXISTS trg_club_ratings_del ON club_ratings;

-- Drop the refresh function
DROP FUNCTION IF EXISTS refresh_club_rating_agg();

COMMIT;