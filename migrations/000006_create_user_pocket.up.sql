CREATE TABLE "user_pocket" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" uuid,
  "pocket_id" uuid
);

ALTER TABLE "user_pocket" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id");
ALTER TABLE "user_pocket" ADD FOREIGN KEY ("pocket_id") REFERENCES "pockets" ("id");