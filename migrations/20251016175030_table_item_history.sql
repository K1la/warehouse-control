-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS item_history (
    id          SERIAL PRIMARY KEY,
    item_id     INTEGER REFERENCES items(id),
    action      VARCHAR(20) NOT NULL,
    old_value   JSONB,
    new_value   JSONB,
    user_id     INTEGER REFERENCES users(id),
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
    );

CREATE INDEX idx_item_history_item_id ON item_history (item_id);
CREATE INDEX idx_item_history_created_at ON item_history (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS item_history;
DROP INDEX idx_item_history_item_id;
DROP INDEX idx_item_history_changed_at;
-- +goose StatementEnd
