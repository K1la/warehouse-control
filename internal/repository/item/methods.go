package item

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/K1la/warehouse-control/internal/model"
)

var (
	ErrItemNotFound = errors.New("no item found")
)

// SetCurrentUser sets the current user in the PostgreSQL session for auditing.
func (r *Postgres) SetCurrentUser(ctx context.Context, userID int64) error {
	_, err := r.db.ExecContext(ctx, "SELECT set_config('app.current_user_id', $1, false)", userID)
	if err != nil {
		return fmt.Errorf("failed to set current_user_id: %w", err)
	}

	return nil
}

func (r *Postgres) Create(ctx context.Context, item model.Item) (int64, error) {
	query := `
	INSERT INTO items
	(name, description, quantity)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		item.Name,
		item.Description,
		item.Quantity,
	).Scan(&item.ID)
	if err != nil {
		return -1, fmt.Errorf("failed to create item: %w", err)
	}

	return item.ID, nil
}

func (r *Postgres) GetAll(ctx context.Context) ([]*model.Item, error) {
	query := `
	SELECT id, name, description, quantity, created_at, updated_at
	FROM items
	ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.Item
	for rows.Next() {
		var it model.Item
		if err = rows.Scan(
			&it.ID,
			&it.Name,
			&it.Description,
			&it.Quantity,
			&it.CreatedAt,
			&it.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, &it)
	}
	return items, nil
}

func (r *Postgres) GetByID(ctx context.Context, itemID int64) (*model.Item, error) {
	query := `
	SELECT id, name, description, quantity, created_at, updated_at
	FROM items
	WHERE id = $1
	`

	var it model.Item
	if err := r.db.QueryRowContext(ctx, query, itemID).Scan(
		&it.ID,
		&it.Name,
		&it.Description,
		&it.Quantity,
		&it.CreatedAt,
		&it.UpdatedAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrItemNotFound
		}
		return nil, err
	}
	return &it, nil
}

func (r *Postgres) Update(ctx context.Context, item model.Item) error {
	query := `
	UPDATE items
	SET name = $1,
	    description = $2,
	    quantity = $3,
	    updated_at = NOW()
	WHERE id = $4
	`
	result, err := r.db.ExecContext(ctx, query,
		item.Name,
		item.Description,
		item.Quantity,
		item.ID,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrItemNotFound
	}

	return nil
}

func (r *Postgres) Delete(ctx context.Context, itemID int64) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM items WHERE id = $1`, itemID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrItemNotFound
	}

	return nil
}
