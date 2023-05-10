-- +goose Up

CREATE TABLE IF NOT EXISTS active_chats (
  tg_chat BIGINT,
  joined_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE active_chats;