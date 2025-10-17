-- +goose Up
-- +goose StatementBegin
CREATE OR REPLACE FUNCTION log_item_insert() RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO item_history(item_id, action, old_value, new_value, user_id, created_at)
    VALUES (
        NEW.id,
        'INSERT',
        NULL,
        to_jsonb(NEW),
        current_setting('app.current_user_id')::INTEGER,
        NOW()
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION log_item_update() RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO item_history(item_id, action, old_value, new_value, user_id, created_at)
    VALUES (
        NEW.id,
        'UPDATE',
        to_jsonb(OLD),
        to_jsonb(NEW),
        current_setting('app.current_user_id')::INTEGER,
        NOW()
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION log_item_delete() RETURNS TRIGGER AS
$$
BEGIN
    INSERT INTO item_history(item_id, action, old_value, user_id, created_at)
    VALUES (
        OLD.id,
        'DELETE',
        to_jsonb(OLD),
        current_setting('app.current_user_id')::INTEGER,
        NOW()
    );
    RETURN OLD;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_item_insert ON items;
CREATE TRIGGER trg_item_insert
    AFTER INSERT ON items
    FOR EACH ROW EXECUTE FUNCTION log_item_insert();

DROP TRIGGER IF EXISTS trg_item_update ON items;
CREATE TRIGGER trg_item_update
    AFTER UPDATE ON items
    FOR EACH ROW EXECUTE FUNCTION log_item_update();

DROP TRIGGER IF EXISTS trg_item_delete ON items;
CREATE TRIGGER trg_item_delete
    BEFORE DELETE ON items
    FOR EACH ROW EXECUTE FUNCTION log_item_delete();

-- CREATE OR REPLACE FUNCTION log_item_changes() RETURNS TRIGGER AS $$
-- BEGIN
--     IF TG_OP = 'INSERT' THEN
--         INSERT INTO history (item_id, action, new_value, user_id)
--         VALUES (NEW.id, TG_OP, row_to_json(NEW), NEW.updated_by);
-- RETURN NEW;
-- ELSIF TG_OP = 'UPDATE' THEN
--         INSERT INTO history (item_id, action, old_value, new_value, user_id)
--         VALUES (NEW.id, TG_OP, row_to_json(OLD), row_to_json(NEW), NEW.updated_by);
-- RETURN NEW;
-- ELSIF TG_OP = 'DELETE' THEN
--         INSERT INTO history (item_id, action, old_value, user_id)
--         VALUES (OLD.id, TG_OP, row_to_json(OLD), OLD.updated_by);
-- RETURN OLD;
-- END IF;
-- RETURN NULL;
-- END;
-- $$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_item_insert ON items;
DROP TRIGGER IF EXISTS trg_item_update ON items;
DROP TRIGGER IF EXISTS trg_item_delete ON items;

DROP FUNCTION IF EXISTS log_item_insert();
DROP FUNCTION IF EXISTS log_item_update();
DROP FUNCTION IF EXISTS log_item_delete();
-- +goose StatementEnd
