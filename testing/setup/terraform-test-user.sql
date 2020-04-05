CREATE SCHEMA terraform_remote_state;
CREATE ROLE terraform_app;
GRANT USAGE,CREATE ON SCHEMA terraform_remote_state to terraform_app;
ALTER DEFAULT PRIVILEGES FOR ROLE terraform_app GRANT INSERT, UPDATE, DELETE, TRUNCATE ON TABLES TO terraform_app;

CREATE USER terraform WITH PASSWORD 'tf';
GRANT CONNECT ON DATABASE terraform TO terraform;
GRANT terraform_app TO terraform;
