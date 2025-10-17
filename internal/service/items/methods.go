package items

import (
	"context"
	"fmt"
	"github.com/K1la/warehouse-control/internal/dto"
	"github.com/K1la/warehouse-control/internal/model"
)

func (s *Service) GetAll(ctx context.Context) ([]*model.Item, error) {
	items, err := s.db.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func (s *Service) GetByID(ctx context.Context, itemID int64) (*model.Item, error) {
	item, err := s.db.GetByID(ctx, itemID)
	if err != nil {
		return nil, err
	}

	return item, nil
}

func (s *Service) Create(ctx context.Context, userID int64, req dto.CreateItemRequest) (int64, error) {
	// Set current user.
	if err := s.db.SetCurrentUser(ctx, userID); err != nil {
		return -1, fmt.Errorf("set current user: %w", err)
	}

	item := model.Item{
		Name:        req.Name,
		Description: req.Description,
		Quantity:    int64(req.Quantity),
	}

	id, err := s.db.Create(ctx, item)
	if err != nil {
		return -1, err
	}

	return id, nil
}

func (s *Service) Update(ctx context.Context, userID, itemID int64, req dto.UpdateItemRequest) error {
	// Set current user.
	if err := s.db.SetCurrentUser(ctx, userID); err != nil {
		return fmt.Errorf("set current user: %w", err)
	}

	item := model.Item{
		ID:          itemID,
		Name:        req.Name,
		Description: req.Description,
		Quantity:    int64(req.Quantity),
	}

	if err := s.db.Update(ctx, item); err != nil {
		return err
	}

	return nil
}

func (s *Service) Delete(ctx context.Context, userID, itemID int64) error {
	// Set current user.
	if err := s.db.SetCurrentUser(ctx, userID); err != nil {
		return fmt.Errorf("set current user: %w", err)
	}

	return s.db.Delete(ctx, itemID)
}
