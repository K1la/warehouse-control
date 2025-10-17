package item_history

import (
	"context"

	"github.com/K1la/warehouse-control/internal/model"
)

func (r *Postgres) ListByID(ctx context.Context, itemID int64) ([]*model.ItemHistory, error) {
	query := `
    SELECT id, COALESCE(item_id, 0) AS item_id, action,
	       COALESCE(old_value::text, '') AS old_value,
	       COALESCE(new_value::text, '') AS new_value,
	       user_id, created_at
	FROM item_history
	WHERE item_id = $1
	ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, itemID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.ItemHistory
	for rows.Next() {
		var it model.ItemHistory
		if err = rows.Scan(
			&it.ID,
			&it.ItemID,
			&it.Action,
			&it.OldValue,
			&it.NewValue,
			&it.UserID,
			&it.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &it)
	}
	return items, nil
}

func (r *Postgres) ListAll(ctx context.Context) ([]*model.ItemHistory, error) {
	query := `
    SELECT id, COALESCE(item_id, 0) AS item_id, action,
	       COALESCE(old_value::text, '') AS old_value,
	       COALESCE(new_value::text, '') AS new_value,
	       user_id, created_at
	FROM item_history
	ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.ItemHistory
	for rows.Next() {
		var it model.ItemHistory
		if err = rows.Scan(
			&it.ID,
			&it.ItemID,
			&it.Action,
			&it.OldValue,
			&it.NewValue,
			&it.UserID,
			&it.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &it)
	}
	return items, nil
}
