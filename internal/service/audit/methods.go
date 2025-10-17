package audit

import (
	"context"

	"github.com/K1la/warehouse-control/internal/model"
)

func (s *Service) ListAll(ctx context.Context) ([]*model.ItemHistory, error) {
	items, err := s.db.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Service) ListByID(ctx context.Context, itemID int64) ([]*model.ItemHistory, error) {
	item, err := s.db.ListByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	return item, nil
}
