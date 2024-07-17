-- Modify "users" table
ALTER TABLE "users" ADD COLUMN "last_login_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP;
