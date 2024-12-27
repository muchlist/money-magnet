ALTER TABLE "users"
  ALTER COLUMN "fcm" DROP DEFAULT;

ALTER TABLE "users"
  ALTER COLUMN "fcm" TYPE varchar(180)[] USING string_to_array("fcm", ',');

ALTER TABLE "users"
  ALTER COLUMN "fcm" SET DEFAULT '{}';