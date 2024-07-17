-- Create "users" table
CREATE TABLE "users" ("id" character varying(255) NOT NULL, "email" character varying(255) NOT NULL, "session" jsonb NOT NULL DEFAULT '{}', "credentials" jsonb NOT NULL DEFAULT '[]', "verified" boolean NOT NULL DEFAULT false, PRIMARY KEY ("id"));
-- Create index "users_email_idx" to table: "users"
CREATE UNIQUE INDEX "users_email_idx" ON "users" ("email");
