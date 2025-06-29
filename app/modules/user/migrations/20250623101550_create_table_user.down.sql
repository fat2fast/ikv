-- Rollback: create_table_user
-- Created at: 2025-06-23 10:15:50

-- Write your down migration here
DROP TABLE IF EXISTS user_users;
DROP TYPE IF EXISTS user_status_enum;
DROP TYPE IF EXISTS user_type_enum;
DROP TYPE IF EXISTS user_role_enum;