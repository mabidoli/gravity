-- Gravity V2 BFF: Initial Database Schema
-- This migration creates the core tables for the unified priority stream

-- Enable UUID generation
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- Users Table
-- Stores participant information for priority items
-- ============================================================================
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE,
    avatar_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Index for email lookups
CREATE INDEX idx_users_email ON users (email);

-- ============================================================================
-- Priority Items Table
-- Core table for the unified stream - represents emails, calendar events, etc.
-- ============================================================================
CREATE TABLE priority_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    source VARCHAR(50) NOT NULL,  -- 'email', 'whatsapp', 'slack', 'teams', 'calendar', 'task', 'youtube', 'linkedin', 'twitter'
    priority VARCHAR(50) NOT NULL, -- 'high', 'medium', 'low'
    is_unread BOOLEAN NOT NULL DEFAULT TRUE,
    snippet TEXT,
    item_timestamp TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Indexes for common query patterns
CREATE INDEX idx_priority_items_user_id ON priority_items (user_id);
CREATE INDEX idx_priority_items_user_timestamp ON priority_items (user_id, item_timestamp DESC);
CREATE INDEX idx_priority_items_user_priority ON priority_items (user_id, priority);
CREATE INDEX idx_priority_items_user_unread ON priority_items (user_id, is_unread) WHERE is_unread = TRUE;

-- ============================================================================
-- Priority Item Participants (Junction Table)
-- Links priority items to their participants
-- ============================================================================
CREATE TABLE priority_item_participants (
    item_id UUID NOT NULL REFERENCES priority_items(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    PRIMARY KEY (item_id, user_id)
);

-- Index for finding all items a user participates in
CREATE INDEX idx_participants_user_id ON priority_item_participants (user_id);

-- ============================================================================
-- Messages Table
-- Stores individual messages within a priority item's conversation
-- ============================================================================
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    item_id UUID NOT NULL REFERENCES priority_items(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(id) ON DELETE SET NULL,
    sender_type VARCHAR(50) NOT NULL, -- 'user', 'other', 'system'
    content_type VARCHAR(50) NOT NULL, -- 'text', 'event', 'social'
    content TEXT,
    full_content_html TEXT,
    message_timestamp TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    -- Complex nested data stored as JSONB for flexibility
    event_details JSONB,      -- CalendarEvent object
    social_details JSONB,     -- SocialContent object
    attachments JSONB,        -- Array of Attachment objects
    ai_insights JSONB         -- Array of AIInsight objects
);

-- Indexes for message queries
CREATE INDEX idx_messages_item_id ON messages (item_id);
CREATE INDEX idx_messages_item_timestamp ON messages (item_id, message_timestamp DESC);
CREATE INDEX idx_messages_sender_id ON messages (sender_id);

-- ============================================================================
-- Trigger for updated_at
-- Automatically updates the updated_at timestamp on row modification
-- ============================================================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_users_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_priority_items_updated_at
    BEFORE UPDATE ON priority_items
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();
