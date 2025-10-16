package dto

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
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity"`
}
