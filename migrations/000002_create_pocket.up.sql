CREATE TABLE IF NOT EXISTS "pockets" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "owner_id" uuid NULL,
  "editor_id" uuid[],
  "watcher_id" uuid[],
  "pocket_name" varchar(100) NOT NULL,
  "icon" int NOT NULL DEFAULT 0,
  "level" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  "version" integer NOT NULL DEFAULT 1
);

ALTER TABLE "pockets" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id") ON DELETE SET NULL;