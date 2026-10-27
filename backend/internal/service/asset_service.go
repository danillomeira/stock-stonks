package service

import (
	"context"

	"github.com/danillomeira/stock-stonks/backend/internal/domain"
	"github.com/danillomeira/stock-stonks/backend/internal/repository"
)

type AssetService struct {
	assets repository.AssetRepository
}

func NewAssetService(assets repository.AssetRepository) *AssetService {
	return &AssetService{assets: assets}
}

func (s *AssetService) Create(ctx context.Context, userID uint, asset *domain.Asset) error {
	asset.UserID = userID
	if err := asset.Validate(); err != nil {
		return err
	}
	return s.assets.Create(ctx, asset)
}

func (s *AssetService) Get(ctx context.Context, userID, id uint) (*domain.Asset, error) {
	return s.assets.GetByID(ctx, userID, id)
}

func (s *AssetService) GetByTicker(ctx context.Context, userID uint, ticker string) (*domain.Asset, error) {
	return s.assets.GetByTicker(ctx, userID, ticker)
}

func (s *AssetService) List(ctx context.Context, userID uint) ([]domain.Asset, error) {
	return s.assets.List(ctx, userID)
}

func (s *AssetService) Update(ctx context.Context, userID uint, asset *domain.Asset) error {
	asset.UserID = userID
	if err := asset.Validate(); err != nil {
		return err
	}
	return s.assets.Update(ctx, asset)
}

func (s *AssetService) Delete(ctx context.Context, userID, id uint) error {
	return s.assets.Delete(ctx, userID, id)
}
