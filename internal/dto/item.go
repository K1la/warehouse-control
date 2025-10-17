package dto

type CreateItemRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity" binding:"required"`
}

type UpdateItemRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
	Quantity    int    `json:"quantity" binding:"required,min=0"`
}
