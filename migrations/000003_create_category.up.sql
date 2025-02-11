CREATE TABLE IF NOT EXISTS "categories" (
  "id" varchar(26) NOT NULL PRIMARY KEY, -- ULID stored as varchar
  "pocket_id" varchar(26) NULL, -- ULID stored as varchar
  "category_name" varchar(100) NOT NULL,
  "is_income" boolean NOT NULL DEFAULT false,
  "default_spend_type" int NOT NULL DEFAULT 0, -- 0:unknown, 1:need, 2:want, 3:saving
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "categories" ADD FOREIGN KEY ("pocket_id") REFERENCES "pockets" ("id") ON DELETE SET NULL;

CREATE UNIQUE INDEX IF NOT EXISTS "unique_category_name_pocket_id"
    ON "categories" ("category_name", "pocket_id");