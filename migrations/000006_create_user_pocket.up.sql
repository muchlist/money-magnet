CREATE TABLE "user_pocket" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" ulid,
  "pocket_id" ulid
);

ALTER TABLE "user_pocket" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE CASCADE;
ALTER TABLE "user_pocket" ADD FOREIGN KEY ("pocket_id") REFERENCES "pockets" ("id") ON DELETE CASCADE;

CREATE UNIQUE INDEX "user_pocket_id" ON "user_pocket" ("user_id", "pocket_id");