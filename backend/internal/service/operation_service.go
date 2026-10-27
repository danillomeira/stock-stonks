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

func (s *OperationService) Create(ctx context.Context, userID uint, op *domain.Operation) error {
	op.UserID = userID
	if err := op.Validate(); err != nil {
		return err
	}
	// GetByID is itself scoped by userID, so this also rejects attaching an
	// operation to an asset that exists but belongs to someone else.
	if _, err := s.assets.GetByID(ctx, userID, op.AssetID); err != nil {
		return fmt.Errorf("operation references unknown asset %d: %w", op.AssetID, err)
	}
	return s.operations.Create(ctx, op)
}

func (s *OperationService) Get(ctx context.Context, userID, id uint) (*domain.Operation, error) {
	return s.operations.GetByID(ctx, userID, id)
}

func (s *OperationService) List(ctx context.Context, filter repository.OperationFilter) ([]domain.Operation, error) {
	return s.operations.List(ctx, filter)
}

func (s *OperationService) Update(ctx context.Context, userID uint, op *domain.Operation) error {
	op.UserID = userID
	if err := op.Validate(); err != nil {
		return err
	}
	return s.operations.Update(ctx, op)
}

func (s *OperationService) Delete(ctx context.Context, userID, id uint) error {
	return s.operations.Delete(ctx, userID, id)
}
