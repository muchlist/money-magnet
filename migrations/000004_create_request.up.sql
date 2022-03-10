CREATE TABLE IF NOT EXISTS "requests" (
  "id" BIGSERIAL PRIMARY KEY,
  "requester" uuid,
  "safe" bigint,
  "safe_name" varchar(100) NOT NULL,
  "is_approved" boolean NOT NULL DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "requests" ADD FOREIGN KEY ("requester") REFERENCES "users" ("id");

ALTER TABLE "requests" ADD FOREIGN KEY ("safe") REFERENCES "safes" ("id") ON DELETE CASCADE;