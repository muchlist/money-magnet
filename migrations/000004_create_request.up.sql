CREATE TABLE IF NOT EXISTS "requests" (
  "id" BIGSERIAL PRIMARY KEY,
  "requester" uuid,
  "approver" uuid DEFAULT NULL,
  "pocket" uuid,
  "pocket_name" varchar(100) NOT NULL,
  "is_approved" boolean DEFAULT NULL,
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "requests" ADD FOREIGN KEY ("requester") REFERENCES "users" ("id");

ALTER TABLE "requests" ADD FOREIGN KEY ("pocket") REFERENCES "pockets" ("id") ON DELETE CASCADE;