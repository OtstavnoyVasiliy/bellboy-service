-- +goose Up
CREATE TABLE IF NOT EXISTS auth_signs (
  chat_id BIGINT NOT NULL,
  user_id INT NOT NULL,
  auth_at TIMESTAMP NOT NULL
);

-- +goose Down
DROP TABLE auth_signs;