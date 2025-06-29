-- Migration: create_table_user
-- Created at: 2025-06-23 10:15:50

-- Write your up migration here

-- Tạo enum types cho PostgreSQL
DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_status_enum') THEN
        CREATE TYPE user_status_enum AS ENUM ('pending','active','inactive','banned','deleted');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_type_enum') THEN
        CREATE TYPE user_type_enum AS ENUM ('email_password','facebook','gmail');
    END IF;
    
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'user_role_enum') THEN
        CREATE TYPE user_role_enum AS ENUM ('user','admin');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS user_users (
    id varchar(36) PRIMARY KEY,
    created_by varchar(36),
    created_at timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by varchar(36),
    updated_at timestamp(6),
    status user_status_enum NOT NULL DEFAULT 'active',
    type user_type_enum NOT NULL DEFAULT 'email_password',
    role user_role_enum NOT NULL DEFAULT 'user',
    first_name varchar(50) NOT NULL,
    last_name varchar(50) NOT NULL,
    phone varchar(20) DEFAULT NULL,
    email varchar(50) NOT NULL,
    password varchar(100) NOT NULL,
    salt varchar(50) NOT NULL,
    UNIQUE (email)
);