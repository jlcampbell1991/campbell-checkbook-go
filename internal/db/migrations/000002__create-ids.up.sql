-- 1. Add column first (nullable for backfill)
ALTER TABLE videos
ADD COLUMN IF NOT EXISTS id UUID;

-- 2. Backfill existing rows (SAFE fallback approach)
-- If uuid_generate_v4() is NOT available, replace this step in Node instead
UPDATE videos
SET id = gen_random_uuid()
WHERE id IS NULL;

-- 3. Enforce NOT NULL after backfill
ALTER TABLE videos
ALTER COLUMN id SET NOT NULL;

-- 4. Set default for future inserts
ALTER TABLE videos
ALTER COLUMN id SET DEFAULT gen_random_uuid();