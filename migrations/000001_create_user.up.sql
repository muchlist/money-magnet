CREATE TABLE IF NOT EXISTS "users" (
  "id" uuid PRIMARY KEY,
  "email" varchar(255) NOT NULL,
  "name" varchar(255) NOT NULL,
  "password" char(60) NOT NULL,
  "roles" text[],
  "pocket_roles" text[],
  "fcm" varchar(255) NOT NULL DEFAULT '',
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

CREATE UNIQUE INDEX "users_email" ON "users" ("email");