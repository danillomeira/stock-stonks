// Package sqlite implements the repository package's interfaces on top of
// GORM with SQLite, with every query wrapped in the shared retry/backoff
// policy from internal/database.
package sqlite

import (
	"context"
	"errors"

	"github.com/danillomeira/stock-stonks/backend/internal/config"
	"github.com/danillomeira/stock-stonks/backend/internal/database"
	"github.com/danillomeira/stock-stonks/backend/internal/domain"
	"github.com/danillomeira/stock-stonks/backend/internal/repository"
	"gorm.io/gorm"
)

type AssetRepository struct {
	db    *gorm.DB
	retry config.RetryConfig
}

func NewAssetRepository(db *gorm.DB, retry config.RetryConfig) *AssetRepository {
	return &AssetRepository{db: db, retry: retry}
}

func (r *AssetRepository) Create(ctx context.Context, asset *domain.Asset) error {
	err := database.WithRetry(ctx, r.retry, func() error {
		return r.db.WithContext(ctx).Create(asset).Error
	})
	if isUniqueConstraintErr(err) {
		return repository.ErrDuplicateTicker
	}
	return err
}

func (r *AssetRepository) GetByID(ctx context.Context, id uint) (*domain.Asset, error) {
	var asset domain.Asset
	err := database.WithRetry(ctx, r.retry, func() error {
		return r.db.WithContext(ctx).First(&asset, id).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *AssetRepository) GetByTicker(ctx context.Context, ticker string) (*domain.Asset, error) {
	var asset domain.Asset
	err := database.WithRetry(ctx, r.retry, func() error {
		return r.db.WithContext(ctx).Where("ticker = ?", ticker).First(&asset).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

func (r *AssetRepository) List(ctx context.Context) ([]domain.Asset, error) {
	var assets []domain.Asset
	err := database.WithRetry(ctx, r.retry, func() error {
		return r.db.WithContext(ctx).Order("ticker").Find(&assets).Error
	})
	return assets, err
}

func (r *AssetRepository) Update(ctx context.Context, asset *domain.Asset) error {
	err := database.WithRetry(ctx, r.retry, func() error {
		return r.db.WithContext(ctx).Save(asset).Error
	})
	if isUniqueConstraintErr(err) {
		return repository.ErrDuplicateTicker
	}
	return err
}

func (r *AssetRepository) Delete(ctx context.Context, id uint) error {
	return database.WithRetry(ctx, r.retry, func() error {
		res := r.db.WithContext(ctx).Delete(&domain.Asset{}, id)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return repository.ErrNotFound
		}
		return nil
	})
}

func isUniqueConstraintErr(err error) bool {
	if err == nil {
		return false
	}
	return errors.Is(err, gorm.ErrDuplicatedKey)
}
