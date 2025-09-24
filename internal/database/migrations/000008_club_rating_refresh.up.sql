BEGIN;
CREATE OR REPLACE FUNCTION refresh_club_rating_agg() RETURNS TRIGGER AS $$
BEGIN
  UPDATE clubs c SET
    rating = COALESCE(sub.avg_rating, 0),
    ratings_count = COALESCE(sub.cnt, 0)
  FROM (
    SELECT club_id, AVG(rating)::float AS avg_rating, COUNT(*) AS cnt
    FROM club_ratings
    WHERE club_id = NEW.club_id
    GROUP BY club_id
  ) sub
  WHERE c.id = NEW.club_id;
  RETURN NULL;
END; $$ LANGUAGE plpgsql;

CREATE TRIGGER trg_club_ratings_ins
AFTER INSERT ON club_ratings
FOR EACH ROW EXECUTE FUNCTION refresh_club_rating_agg();

CREATE TRIGGER trg_club_ratings_upd
AFTER UPDATE ON club_ratings
FOR EACH ROW EXECUTE FUNCTION refresh_club_rating_agg();

CREATE TRIGGER trg_club_ratings_del
AFTER DELETE ON club_ratings
FOR EACH ROW EXECUTE FUNCTION refresh_club_rating_agg();
COMMIT;