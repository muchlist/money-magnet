CREATE TABLE IF NOT EXISTS "categories" (
  "id" uuid PRIMARY KEY,
  "pocket_id" uuid NULL,
  "category_name" varchar(100) NOT NULL,
  "is_income" boolean NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "categories" ADD FOREIGN KEY ("pocket_id") REFERENCES "pockets" ("id") ON DELETE SET NULL;