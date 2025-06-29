-- Rollback: create_table_book
-- Created at: 2025-06-28 14:32:15

-- Write your down migration here
DROP TABLE IF EXISTS book_books;
DROP TYPE IF EXISTS book_status_enum;