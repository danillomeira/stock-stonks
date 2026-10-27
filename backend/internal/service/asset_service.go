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

func (s *AssetService) Create(ctx context.Context, asset *domain.Asset) error {
	if err := asset.Validate(); err != nil {
		return err
	}
	return s.assets.Create(ctx, asset)
}

func (s *AssetService) Get(ctx context.Context, id uint) (*domain.Asset, error) {
	return s.assets.GetByID(ctx, id)
}

func (s *AssetService) GetByTicker(ctx context.Context, ticker string) (*domain.Asset, error) {
	return s.assets.GetByTicker(ctx, ticker)
}

func (s *AssetService) List(ctx context.Context) ([]domain.Asset, error) {
	return s.assets.List(ctx)
}

func (s *AssetService) Update(ctx context.Context, asset *domain.Asset) error {
	if err := asset.Validate(); err != nil {
		return err
	}
	return s.assets.Update(ctx, asset)
}

func (s *AssetService) Delete(ctx context.Context, id uint) error {
	return s.assets.Delete(ctx, id)
}
