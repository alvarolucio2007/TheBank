CREATE TABLE "users"(
  "username" VARCHAR PRIMARY KEY,
  "hashed_password" VARCHAR NOT NULL,
  "full_name" VARCHAR NOT NULL,
  "email" VARCHAR UNIQUE NOT NULL,
  "password_changed_at" TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);
CREATE TABLE "accounts"(
  "id" BIGSERIAL PRIMARY KEY,
  "owner" VARCHAR NOT NULL REFERENCES "users"("username"),
  "balance" BIGINT NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);
CREATE TABLE "entries"(
  "id" BIGSERIAL PRIMARY KEY,
  "account_id" BIGINT NOT NULL REFERENCES "accounts"("id"),
  "amount" BIGINT NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);
CREATE TABLE "transfers"(
  "id" BIGSERIAL PRIMARY KEY,
  "from_account_id" BIGINT NOT NULL REFERENCES "accounts"("id"),
  "to_account_id" BIGINT NOT NULL REFERENCES "accounts"("id"),
  "amount" BIGINT NOT NULL,
  "created_at"  TIMESTAMPTZ NOT NULL DEFAULT (now())
);
CREATE INDEX ON "accounts" ("owner");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");


COMMENT ON COLUMN "entries"."amount" IS 'can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'must be positive';
