package model

import (
	"time"
)

type ItemHistory struct {
	ID        int64     `db:"id" json:"id"`
	ItemID    int64     `db:"item_id" json:"item_id"`
	Action    string    `db:"action" json:"action"`
	OldValue  string    `db:"old_value" json:"old_value"`
	NewValue  string    `db:"new_value" json:"new_value"`
	UserID    int64     `db:"user_id" json:"user_id"`
	CreatedAt time.Time `db:"created_at" json:"created_at"`
}
