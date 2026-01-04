-- Migration: Change user_id to support Clerk string IDs
-- Clerk uses string IDs like 'user_37nulvJh1uJOvVawhJvG5IgQqiR' instead of UUIDs

-- Add clerk_id column to users table
ALTER TABLE users ADD COLUMN IF NOT EXISTS clerk_id VARCHAR(255) UNIQUE;
CREATE INDEX IF NOT EXISTS idx_users_clerk_id ON users (clerk_id);

-- Alter priority_items to use TEXT for user_id (to store Clerk IDs directly)
-- First, drop the foreign key constraint
ALTER TABLE priority_items DROP CONSTRAINT IF EXISTS priority_items_user_id_fkey;

-- Change the column type
ALTER TABLE priority_items ALTER COLUMN user_id TYPE VARCHAR(255) USING user_id::VARCHAR(255);

-- Update priority_item_participants similarly
ALTER TABLE priority_item_participants DROP CONSTRAINT IF EXISTS priority_item_participants_user_id_fkey;
ALTER TABLE priority_item_participants ALTER COLUMN user_id TYPE VARCHAR(255) USING user_id::VARCHAR(255);

-- Note: We're removing the FK constraint for now.
-- In production, you'd want to create a proper user sync flow.
