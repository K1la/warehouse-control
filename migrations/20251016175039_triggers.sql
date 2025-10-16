-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION log_item_changes() RETURNS TRIGGER AS $$
BEGIN
    IF TG_OP = 'INSERT' THEN
        INSERT INTO history (item_id, action, new_value, user_id)
        VALUES (NEW.id, TG_OP, row_to_json(NEW), NEW.updated_by);
RETURN NEW;
ELSIF TG_OP = 'UPDATE' THEN
        INSERT INTO history (item_id, action, old_value, new_value, user_id)
        VALUES (NEW.id, TG_OP, row_to_json(OLD), row_to_json(NEW), NEW.updated_by);
RETURN NEW;
ELSIF TG_OP = 'DELETE' THEN
        INSERT INTO history (item_id, action, old_value, user_id)
        VALUES (OLD.id, TG_OP, row_to_json(OLD), OLD.updated_by);
RETURN OLD;
END IF;
RETURN NULL;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_item_log ON items;
CREATE TRIGGER trg_item_log
AFTER INSERT OR UPDATE OR DELETE ON items
FOR EACH ROW EXECUTE FUNCTION log_item_changes();
-- +goose StatementEnd
