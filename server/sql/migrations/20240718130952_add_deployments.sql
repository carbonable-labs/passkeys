-- Create "accounts" table
CREATE TABLE "accounts" ("id" character varying(255) NOT NULL, "user_id" character varying(255) NOT NULL, "address" character varying(255) NULL DEFAULT NULL::character varying, "data" jsonb NOT NULL DEFAULT '{}', "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("id"), CONSTRAINT "accounts_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "accounts_address_idx" to table: "accounts"
CREATE UNIQUE INDEX "accounts_address_idx" ON "accounts" ("user_id", "address");
-- Create "account_deployments" table
CREATE TABLE "account_deployments" ("id" character varying(255) NOT NULL, "user_id" character varying(255) NOT NULL, "account_id" character varying(255) NOT NULL, "status" character varying(255) NOT NULL, "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("id"), CONSTRAINT "account_deployments_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "account_deployments_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "accounts_deployments_idx" to table: "account_deployments"
CREATE UNIQUE INDEX "accounts_deployments_idx" ON "account_deployments" ("user_id", "account_id");
