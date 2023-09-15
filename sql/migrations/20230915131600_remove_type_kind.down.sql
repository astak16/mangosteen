BEGIN;
CREATE TYPE kind AS ENUM ('expenses', 'in_come', '');
ALTER TABLE items DROP COLUMN kind;
ALTER TABLE items ADD COLUMN kind kind NOT NULL DEFAULT 'expenses';
ALTER TABLE tags DROP COLUMN kind;
ALTER TABLE tags ADD COLUMN kind kind NOT NULL DEFAULT 'expenses';
COMMIT;
