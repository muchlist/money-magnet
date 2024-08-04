CREATE TABLE IF NOT EXISTS "requests" (
  "id" BIGSERIAL PRIMARY KEY,
  "requester_id" ulid,
  "approver_id" ulid DEFAULT NULL,
  "pocket_id" ulid,
  "pocket_name" varchar(100) NOT NULL,
  "is_approved" boolean DEFAULT false,
  "is_rejected" boolean DEFAULT false,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "requests" ADD FOREIGN KEY ("requester_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "requests" ADD FOREIGN KEY ("approver_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "requests" ADD FOREIGN KEY ("pocket_id") REFERENCES "pockets" ("id") ON DELETE CASCADE;