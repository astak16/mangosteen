BEGIN;
ALTER TABLE items
ALTER COLUMN kind TYPE VARCHAR(100),
ALTER COLUMN kind SET DEFAULT 'expenses';

ALTER TABLE tags
ALTER COLUMN kind TYPE VARCHAR(100),
ALTER COLUMN kind SET DEFAULT 'expenses';
DROP TYPE kind;
COMMIT;
