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

type OperationRepository struct {
	db    *gorm.DB
	retry config.RetryConfig
}

func NewOperationRepository(db *gorm.DB, retry config.RetryConfig) *OperationRepository {
	return &OperationRepository{db: db, retry: retry}
}

func (r *OperationRepository) Create(ctx context.Context, op *domain.Operation) error {
	return database.WithRetry(ctx, r.retry, func() error {
		return r.db.WithContext(ctx).Create(op).Error
	})
}

func (r *OperationRepository) GetByID(ctx context.Context, id uint) (*domain.Operation, error) {
	var op domain.Operation
	err := database.WithRetry(ctx, r.retry, func() error {
		return r.db.WithContext(ctx).Preload("Asset").First(&op, id).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &op, nil
}

func (r *OperationRepository) List(ctx context.Context, filter repository.OperationFilter) ([]domain.Operation, error) {
	var ops []domain.Operation
	err := database.WithRetry(ctx, r.retry, func() error {
		q := r.db.WithContext(ctx).Preload("Asset").Order("date")
		if filter.AssetID != nil {
			q = q.Where("asset_id = ?", *filter.AssetID)
		}
		if filter.From != nil {
			q = q.Where("date >= ?", *filter.From)
		}
		if filter.To != nil {
			q = q.Where("date <= ?", *filter.To)
		}
		return q.Find(&ops).Error
	})
	return ops, err
}

func (r *OperationRepository) Update(ctx context.Context, op *domain.Operation) error {
	return database.WithRetry(ctx, r.retry, func() error {
		return r.db.WithContext(ctx).Save(op).Error
	})
}

func (r *OperationRepository) Delete(ctx context.Context, id uint) error {
	return database.WithRetry(ctx, r.retry, func() error {
		res := r.db.WithContext(ctx).Delete(&domain.Operation{}, id)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return repository.ErrNotFound
		}
		return nil
	})
}
