package service

import (
	"context"
	"fmt"

	"github.com/danillomeira/stock-stonks/backend/internal/domain"
	"github.com/danillomeira/stock-stonks/backend/internal/repository"
)

type OperationService struct {
	operations repository.OperationRepository
	assets     repository.AssetRepository
}

func NewOperationService(operations repository.OperationRepository, assets repository.AssetRepository) *OperationService {
	return &OperationService{operations: operations, assets: assets}
}

func (s *OperationService) Create(ctx context.Context, op *domain.Operation) error {
	if err := op.Validate(); err != nil {
		return err
	}
	if _, err := s.assets.GetByID(ctx, op.AssetID); err != nil {
		return fmt.Errorf("operation references unknown asset %d: %w", op.AssetID, err)
	}
	return s.operations.Create(ctx, op)
}

func (s *OperationService) Get(ctx context.Context, id uint) (*domain.Operation, error) {
	return s.operations.GetByID(ctx, id)
}

func (s *OperationService) List(ctx context.Context, filter repository.OperationFilter) ([]domain.Operation, error) {
	return s.operations.List(ctx, filter)
}

func (s *OperationService) Update(ctx context.Context, op *domain.Operation) error {
	if err := op.Validate(); err != nil {
		return err
	}
	return s.operations.Update(ctx, op)
}

func (s *OperationService) Delete(ctx context.Context, id uint) error {
	return s.operations.Delete(ctx, id)
}
