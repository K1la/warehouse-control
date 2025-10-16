package model

import "time"

type Item struct {
	ID          int64
	Name        string
	Description string
	Quantity    int64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
