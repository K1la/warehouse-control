package items

import (
	"context"
	"github.com/K1la/warehouse-control/internal/dto"
	"github.com/K1la/warehouse-control/internal/model"
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
	Create(ctx context.Context, userID int64, req dto.CreateItemRequest) (int64, error)
	Update(ctx context.Context, userID, itemID int64, req dto.UpdateItemRequest) error
	Delete(ctx context.Context, userID, itemID int64) error
	GetByID(ctx context.Context, itemID int64) (*model.Item, error)
	GetAll(ctx context.Context) ([]*model.Item, error)
}
