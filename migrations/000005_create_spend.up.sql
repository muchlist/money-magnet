CREATE TABLE  IF NOT EXISTS "spends" (
  "id" uuid PRIMARY KEY,
  "user" uuid,
  "pocket" uuid,
  "category" uuid,
  "category_name" varchar(100),
  "category_x" uuid,
  "category_name_x" varchar(100),
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

ALTER TABLE "spends" ADD FOREIGN KEY ("user") REFERENCES "users" ("id");

ALTER TABLE "spends" ADD FOREIGN KEY ("pocket") REFERENCES "pockets" ("id") ON DELETE CASCADE;

ALTER TABLE "spends" ADD FOREIGN KEY ("category") REFERENCES "categories" ("id");

ALTER TABLE "spends" ADD FOREIGN KEY ("category_x") REFERENCES "categories" ("id");

CREATE INDEX "spend_pocket_date" ON "spends" ("pocket", "date");

CREATE INDEX "spend_pocket_user" ON "spends" ("pocket", "user");