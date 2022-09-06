CREATE TABLE  IF NOT EXISTS "spends" (
  "id" uuid DEFAULT gen_random_uuid() PRIMARY KEY,
  "user_id" uuid NULL,
  "pocket_id" uuid NULL,
  "category_id" uuid NULL,
  "category_id_2" uuid NULL,
  "name" varchar(255) NOT NULL,
  "price" bigint NOT NULL,
  "balance" bigint NOT NULL,
  "is_income" boolean NOT NULL DEFAULT false,
  "type" int NOT NULL,
  "date" timestamp NOT NULL DEFAULT (now()),
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now()),
  "version" integer NOT NULL DEFAULT 1
);

ALTER TABLE "spends" ADD FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON DELETE SET NULL;
ALTER TABLE "spends" ADD FOREIGN KEY ("pocket_id") REFERENCES "pockets" ("id") ON DELETE SET NULL;
ALTER TABLE "spends" ADD FOREIGN KEY ("category_id") REFERENCES "categories" ("id") ON DELETE SET NULL;
ALTER TABLE "spends" ADD FOREIGN KEY ("category_id_2") REFERENCES "categories" ("id") ON DELETE SET NULL;

CREATE INDEX "spend_pocket_date" ON "spends" ("pocket_id", "date");
CREATE INDEX "spend_pocket_user" ON "spends" ("pocket_id", "user_id");