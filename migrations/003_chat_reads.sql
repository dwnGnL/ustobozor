ALTER TABLE chat_messages
    ADD COLUMN IF NOT EXISTS read_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_chat_messages_read_at ON chat_messages(read_at);
