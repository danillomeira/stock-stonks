package service

import (
	"context"

	"github.com/danillomeira/stock-stonks/backend/internal/domain"
	"github.com/danillomeira/stock-stonks/backend/internal/repository"
)

// PositionService derives current holdings from the full operation history.
// This is the groundwork the future IR export feature will build on: the
// export needs exactly this average-cost replay, per asset, per tax year.
type PositionService struct {
	operations repository.OperationRepository
}

func NewPositionService(operations repository.OperationRepository) *PositionService {
	return &PositionService{operations: operations}
}

// ForAsset replays every Operation for assetID in date order and returns the
// resulting Position.
func (s *PositionService) ForAsset(ctx context.Context, assetID uint) (domain.Position, error) {
	filter := repository.OperationFilter{AssetID: &assetID}
	ops, err := s.operations.List(ctx, filter)
	if err != nil {
		return domain.Position{}, err
	}

	pos := domain.Position{AssetID: assetID}
	for _, op := range ops {
		pos.ApplyOperation(op)
	}
	return pos, nil
}

// All returns the current Position for every asset that has at least one operation.
func (s *PositionService) All(ctx context.Context) ([]domain.Position, error) {
	ops, err := s.operations.List(ctx, repository.OperationFilter{})
	if err != nil {
		return nil, err
	}

	order := make([]uint, 0)
	byAsset := make(map[uint]*domain.Position)
	for _, op := range ops {
		pos, ok := byAsset[op.AssetID]
		if !ok {
			pos = &domain.Position{AssetID: op.AssetID}
			byAsset[op.AssetID] = pos
			order = append(order, op.AssetID)
		}
		pos.ApplyOperation(op)
	}

	result := make([]domain.Position, 0, len(order))
	for _, assetID := range order {
		result = append(result, *byAsset[assetID])
	}
	return result, nil
}
