CREATE TABLE IF NOT EXISTS users (
  id VARCHAR(255) PRIMARY KEY,
  email VARCHAR(255) NOT NULL,
  session JSONB NOT NULL DEFAULT '{}',
  credentials JSONB NOT NULL DEFAULT '[]',
  verified BOOLEAN NOT NULL DEFAULT FALSE,
  last_login_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS users_email_idx ON users (email);

CREATE TABLE IF NOT EXISTS accounts (
  id VARCHAR(255) PRIMARY KEY,
  user_id VARCHAR(255) NOT NULL REFERENCES users (id),
  address VARCHAR(255) DEFAULT NULL,
  data JSONB NOT NULL DEFAULT '{}',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(user_id) REFERENCES users(id)
);
CREATE UNIQUE INDEX IF NOT EXISTS accounts_address_idx ON accounts (user_id, address);

CREATE TABLE IF NOT EXISTS account_deployments (
  id VARCHAR(255) PRIMARY KEY,
  user_id VARCHAR(255) NOT NULL REFERENCES users (id),
  account_id VARCHAR(255) NOT NULL REFERENCES accounts (id),
  status VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(account_id) REFERENCES accounts(id), 
  FOREIGN KEY(user_id) REFERENCES users(id)

);
CREATE UNIQUE INDEX IF NOT EXISTS accounts_deployments_idx ON account_deployments (user_id, account_id);

CREATE TABLE IF NOT EXISTS account_deployment_logs (
  id VARCHAR(255) PRIMARY KEY,
  account_id VARCHAR(255) NOT NULL REFERENCES accounts (id),
  message VARCHAR(255) NOT NULL,
  payload JSONB NOT NULL DEFAULT '{}',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  FOREIGN KEY(account_id) REFERENCES accounts(id)
);
CREATE UNIQUE INDEX IF NOT EXISTS account_deployment_logs_idx ON account_deployment_logs (account_id);

