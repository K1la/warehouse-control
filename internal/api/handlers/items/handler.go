package items

import (
	"context"
	"github.com/K1la/warehouse-control/internal/dto"
	"github.com/rs/zerolog"
)

type Handler struct {
	service Service
	log     zerolog.Logger
}

func New(s Service, l zerolog.Logger) *Handler {
	return &Handler{service: s, log: l}
}

type Service interface {
	Create(ctx context.Context, req dto.CreateItemRequest) (*dto.ItemResponse, error)
	GetAll(ctx context.Context) ([]dto.ItemResponse, error)
	GetByID(ctx context.Context, id string) (*dto.ItemResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateItemRequest) (*dto.ItemResponse, error)
	Delete(ctx context.Context, id string) error
}
