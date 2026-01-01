-- Rollback: Drop all tables and extensions created in the initial schema

-- Drop triggers first
DROP TRIGGER IF EXISTS update_priority_items_updated_at ON priority_items;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop the trigger function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order of creation (respecting foreign key constraints)
DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS priority_item_participants;
DROP TABLE IF EXISTS priority_items;
DROP TABLE IF EXISTS users;

-- Drop extensions
DROP EXTENSION IF EXISTS "uuid-ossp";
