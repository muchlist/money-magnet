CREATE TABLE  IF NOT EXISTS "spends" (
  "id" uuid PRIMARY KEY,
  "user" uuid,
  "safe" bigint,
  "category" uuid,
  "category_name" varchar(100),
  "name" varchar(255) NOT NULL,
  "price" bigint NOT NULL,
  "balance" bigint NOT NULL,
  "is_income" boolean NOT NULL DEFAULT false,
  "type" int NOT NULL,
  "date" timestamp NOT NULL DEFAULT (now()),
  "created_at" timestamp NOT NULL DEFAULT (now()),
  "updated_at" timestamp NOT NULL DEFAULT (now())
);

ALTER TABLE "spends" ADD FOREIGN KEY ("user") REFERENCES "users" ("id");

ALTER TABLE "spends" ADD FOREIGN KEY ("safe") REFERENCES "safes" ("id") ON DELETE CASCADE;

ALTER TABLE "spends" ADD FOREIGN KEY ("category") REFERENCES "categories" ("id");

CREATE INDEX "spend_safe_date" ON "spends" ("safe", "date");

CREATE INDEX "spend_safe_user" ON "spends" ("safe", "user");