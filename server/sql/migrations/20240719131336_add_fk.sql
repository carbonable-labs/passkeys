-- Modify "accounts" table
ALTER TABLE "accounts" ADD CONSTRAINT "accounts_user_id_fkey1" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
-- Create "account_deployment_logs" table
CREATE TABLE "account_deployment_logs" ("id" character varying(255) NOT NULL, "account_id" character varying(255) NOT NULL, "message" character varying(255) NOT NULL, "payload" jsonb NOT NULL DEFAULT '{}', "created_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, "updated_at" timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY ("id"), CONSTRAINT "account_deployment_logs_account_id_fkey" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, CONSTRAINT "account_deployment_logs_account_id_fkey1" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION);
-- Create index "account_deployment_logs_idx" to table: "account_deployment_logs"
CREATE UNIQUE INDEX "account_deployment_logs_idx" ON "account_deployment_logs" ("account_id");
-- Modify "account_deployments" table
ALTER TABLE "account_deployments" ADD CONSTRAINT "account_deployments_account_id_fkey1" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION, ADD CONSTRAINT "account_deployments_user_id_fkey1" FOREIGN KEY ("user_id") REFERENCES "users" ("id") ON UPDATE NO ACTION ON DELETE NO ACTION;
