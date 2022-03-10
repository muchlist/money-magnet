CREATE TABLE IF NOT EXISTS "safes" (
  "id" BIGSERIAL PRIMARY KEY,
  "owner" uuid,
  "editor" uuid[],
  "watcher" uuid[],
  "safe_name" varchar(100) NOT NULL,
  "level" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "safes" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id");