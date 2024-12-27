ALTER TABLE "users"
  ALTER COLUMN "fcm" DROP DEFAULT;

ALTER TABLE "users"
  ALTER COLUMN "fcm" TYPE varchar(255) USING array_to_string("fcm", ',');

ALTER TABLE "users"
  ALTER COLUMN "fcm" SET DEFAULT '';