ALTER TABLE IF EXISTS categories
ADD COLUMN IF NOT EXISTS category_icon integer DEFAULT 0 NOT NULL;

ALTER TABLE categories
ALTER COLUMN id SET DEFAULT gen_random_uuid();