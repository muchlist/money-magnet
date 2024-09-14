CREATE TABLE IF NOT EXISTS "pockets" (
  "id" varchar(26) NOT NULL PRIMARY KEY, -- ULID stored as varchar
  "owner_id" varchar(26) NULL,
  "editor_id" text[],
  "watcher_id" text[],
  "pocket_name" varchar(100) NOT NULL,
  "balance" bigint NOT NULL DEFAULT 0,
  "currency" varchar(10) NOT NULL DEFAULT '',
  "icon" int NOT NULL DEFAULT 0,
  "level" int NOT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  "version" integer NOT NULL DEFAULT 1
);

ALTER TABLE "pockets" ADD FOREIGN KEY ("owner_id") REFERENCES "users" ("id") ON DELETE SET NULL;