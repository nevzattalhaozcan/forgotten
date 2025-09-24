BEGIN;

-- allow NULL owner_id
ALTER TABLE clubs ALTER COLUMN owner_id DROP NOT NULL;

-- drop old FK (name may vary by environment; drop both if present)
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.table_constraints
        WHERE constraint_name = 'fk_users_owned_clubs' AND table_name = 'clubs'
    ) THEN
        ALTER TABLE clubs DROP CONSTRAINT fk_users_owned_clubs;
    END IF;
EXCEPTION WHEN undefined_object THEN
    -- ignore
END$$;

-- try dropping any default-named FK as a fallback
DO $$
DECLARE
    conname text;
BEGIN
    SELECT tc.constraint_name INTO conname
    FROM information_schema.table_constraints tc
    JOIN information_schema.key_column_usage kcu
      ON tc.constraint_name = kcu.constraint_name
    WHERE tc.table_name = 'clubs'
      AND tc.constraint_type = 'FOREIGN KEY'
      AND kcu.column_name = 'owner_id'
    LIMIT 1;
    IF conname IS NOT NULL THEN
        EXECUTE format('ALTER TABLE clubs DROP CONSTRAINT %I', conname);
    END IF;
END$$;

-- recreate with SET NULL
ALTER TABLE clubs
  ADD CONSTRAINT fk_users_owned_clubs
  FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE SET NULL;

COMMIT;