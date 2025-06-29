-- Migration: create_table_book
-- Created at: 2025-06-28 14:32:15

-- Write your up migration here

DO $$ 
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'book_status_enum') THEN
        CREATE TYPE book_status_enum AS ENUM ('pending','active','inactive','banned','deleted');
    END IF;
END $$;

CREATE TABLE IF NOT EXISTS book_books (
   id varchar(36) PRIMARY KEY,
    created_by varchar(36),
    created_at timestamp(6) NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_by varchar(36),
    updated_at timestamp(6),
    status book_status_enum NOT NULL DEFAULT 'active',
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    description TEXT,
    price DECIMAL(10, 2) NOT NULL,
    published_at timestamp(6),
    cover_image VARCHAR(255),
    UNIQUE (title)
);