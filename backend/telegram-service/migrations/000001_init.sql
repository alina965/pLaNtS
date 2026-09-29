CREATE TABLE IF NOT EXISTS telegram_links (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL UNIQUE,
    chat_id BIGINT UNIQUE,
    link_code TEXT UNIQUE,
    created_at TIMESTAMPTZ NOT NULL,
    linked_at TIMESTAMPTZ
);