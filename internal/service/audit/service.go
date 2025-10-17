package audit

import (
	"context"

	"github.com/K1la/warehouse-control/internal/model"
	"github.com/rs/zerolog"
)

type Service struct {
	db  Repo
	log zerolog.Logger
}

func New(d Repo, l zerolog.Logger) *Service {
	return &Service{db: d, log: l}
}

type Repo interface {
	ListByID(ctx context.Context, itemID int64) ([]*model.ItemHistory, error)
	ListAll(ctx context.Context) ([]*model.ItemHistory, error)
}
