package items

import (
	"context"
	"github.com/K1la/warehouse-control/internal/dto"
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
	Create(ctx context.Context, req *model.Item) error
	GetAll(ctx context.Context) ([]dto.ItemResponse, error)
	GetByID(ctx context.Context, id string) (*dto.ItemResponse, error)
	Update(ctx context.Context, req *dto.ItemResponse) error
	Delete(ctx context.Context, id string) error
}

type Service interface {
	Create(ctx context.Context, req dto.CreateItemRequest) (*dto.ItemResponse, error)
	GetAll(ctx context.Context) ([]dto.ItemResponse, error)
	GetByID(ctx context.Context, id string) (*dto.ItemResponse, error)
	Update(ctx context.Context, id string, req dto.UpdateItemRequest) (*dto.ItemResponse, error)
	Delete(ctx context.Context, id string) error
}
