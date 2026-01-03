CREATE TABLE IF NOT EXISTS chats (
    id UUID PRIMARY KEY,
    request_id UUID NOT NULL REFERENCES job_requests(id) ON DELETE CASCADE,
    creator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    initiator_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL,
    last_message_at TIMESTAMPTZ,
    CHECK (creator_id <> initiator_id)
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_chats_request_initiator
    ON chats(request_id, initiator_id);
CREATE INDEX IF NOT EXISTS idx_chats_participants
    ON chats(creator_id, initiator_id);
CREATE INDEX IF NOT EXISTS idx_chats_last_message
    ON chats(last_message_at);

CREATE TABLE IF NOT EXISTS chat_messages (
    id UUID PRIMARY KEY,
    chat_id UUID NOT NULL REFERENCES chats(id) ON DELETE CASCADE,
    sender_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    text TEXT,
    photo_path TEXT,
    created_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_chat_messages_chat_id
    ON chat_messages(chat_id);
CREATE INDEX IF NOT EXISTS idx_chat_messages_created_at
    ON chat_messages(created_at);
