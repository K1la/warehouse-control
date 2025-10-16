package items

import (
	"context"
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
	ListByID(ctx context.Context, itemID int64) ([]*model.ItemHistory, error)
	ListAll(ctx context.Context) ([]*model.ItemHistory, error)
}
