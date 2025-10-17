-- +goose Up
-- +goose StatementBegin
ALTER TABLE item_history
    DROP CONSTRAINT IF EXISTS item_history_item_id_fkey;

ALTER TABLE item_history
    ADD CONSTRAINT item_history_item_id_fkey
        FOREIGN KEY (item_id)
        REFERENCES items(id)
        ON DELETE SET NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE item_history
    DROP CONSTRAINT IF EXISTS item_history_item_id_fkey;

-- restore default behavior (no action)
ALTER TABLE item_history
    ADD CONSTRAINT item_history_item_id_fkey
        FOREIGN KEY (item_id)
        REFERENCES items(id);
-- +goose StatementEnd


