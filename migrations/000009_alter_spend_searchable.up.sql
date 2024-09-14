CREATE EXTENSION IF NOT EXISTS pg_trgm;
-- CREATE EXTENSION IF NOT EXISTS fuzzystrmatch;

-- index trigram for search by spend name
CREATE INDEX trgm_spend_name ON "spends" USING gin ("name" gin_trgm_ops);