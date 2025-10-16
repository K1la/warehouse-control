package dto

import "time"

type CreateItemRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity" binding:"required"`
}

type UpdateItemRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Quantity    *int   `json:"quantity" binding:"required,min=0"`
}

type ItemResponse struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Quantity    int       `json:"quantity"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
