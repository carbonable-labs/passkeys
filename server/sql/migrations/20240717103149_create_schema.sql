-- Create "users" table
CREATE TABLE "users" ("id" character varying(255) NOT NULL, "email" character varying(255) NOT NULL, "session" jsonb NOT NULL DEFAULT '{}', PRIMARY KEY ("id"));
