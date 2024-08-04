CREATE EXTENSION IF NOT EXISTS citext;
CREATE EXTENSION IF NOT EXISTS ulid;

CREATE TABLE IF NOT EXISTS "users" (
  "id" ulid DEFAULT gen_ulid() PRIMARY KEY,
  "email" citext NOT NULL,
  "name" varchar(255) NOT NULL,
  "password" bytea NOT NULL,
  "roles" text[],
  "pocket_roles" text[],
  "fcm" varchar(255) NOT NULL DEFAULT '',
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  "version" integer NOT NULL DEFAULT 1
);

CREATE UNIQUE INDEX "users_email" ON "users" ("email");
CREATE INDEX "name" ON "users" ("name");