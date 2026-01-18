-- +goose Up
-- SQL section 'Up' is executed when running 'goose up'

CREATE TABLE chats (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(200) NOT NULL CHECK (length(trim(title)) > 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    chat_id BIGINT NOT NULL,
    text VARCHAR(5000) NOT NULL CHECK (length(trim(text)) > 0),
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_chat_messages
        FOREIGN KEY (chat_id)
            REFERENCES chats(id)
            ON DELETE CASCADE
);

CREATE INDEX idx_messages_chat_id ON messages(chat_id);


-- +goose Down
-- SQL section 'Down' is executed when running 'goose down'

DROP TABLE IF EXISTS messages;
DROP TABLE IF EXISTS chats;