CREATE TABLE IF NOT EXISTS "categories" (
  "id" uuid PRIMARY KEY,
  "safe" bigint,
  "category_name" varchar(100) NOT NULL,
  "is_income" boolean NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "categories" ADD FOREIGN KEY ("safe") REFERENCES "safes" ("id") ON DELETE CASCADE;