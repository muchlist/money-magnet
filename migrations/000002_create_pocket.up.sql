CREATE TABLE IF NOT EXISTS "pockets" (
  "id" BIGSERIAL PRIMARY KEY,
  "owner" uuid,
  "editor" uuid[],
  "watcher" uuid[],
  "pocket_name" varchar(100) NOT NULL,
  "icon" int NOT NULL DEFAULT 0,
  "level" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  "version" integer NOT NULL DEFAULT 1
);

ALTER TABLE "pockets" ADD FOREIGN KEY ("owner") REFERENCES "users" ("id");