package items

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
	SetCurrentUser(ctx context.Context, userID int64) error
	Create(ctx context.Context, item model.Item) (int64, error)
	GetAll(ctx context.Context) ([]*model.Item, error)
	GetByID(ctx context.Context, itemID int64) (*model.Item, error)
	Update(ctx context.Context, item model.Item) error
	Delete(ctx context.Context, itemID int64) error
}
