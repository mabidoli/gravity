-- Rollback: Revert user_id changes
-- Note: This may fail if there are Clerk IDs that can't be converted to UUIDs

ALTER TABLE priority_items ALTER COLUMN user_id TYPE UUID USING user_id::UUID;
ALTER TABLE priority_item_participants ALTER COLUMN user_id TYPE UUID USING user_id::UUID;

ALTER TABLE priority_items
    ADD CONSTRAINT priority_items_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

ALTER TABLE priority_item_participants
    ADD CONSTRAINT priority_item_participants_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

DROP INDEX IF EXISTS idx_users_clerk_id;
ALTER TABLE users DROP COLUMN IF EXISTS clerk_id;
