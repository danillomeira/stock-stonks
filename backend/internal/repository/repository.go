// Package repository defines the persistence interfaces the service layer
// depends on. Concrete implementations (currently sqlite/) are swapped in at
// wiring time in cmd/api/main.go, so the service layer never imports GORM or
// any driver directly.
package repository

import (
	"context"
	"errors"
	"time"

	"github.com/danillomeira/stock-stonks/backend/internal/domain"
)

// ErrNotFound is returned by any Get/Update/Delete that targets a missing row.
var ErrNotFound = errors.New("repository: entity not found")

// ErrDuplicateTicker is returned when creating an Asset whose ticker already exists.
var ErrDuplicateTicker = errors.New("repository: ticker already exists")

type AssetRepository interface {
	Create(ctx context.Context, asset *domain.Asset) error
	GetByID(ctx context.Context, id uint) (*domain.Asset, error)
	GetByTicker(ctx context.Context, ticker string) (*domain.Asset, error)
	List(ctx context.Context) ([]domain.Asset, error)
	Update(ctx context.Context, asset *domain.Asset) error
	Delete(ctx context.Context, id uint) error
}

// OperationFilter narrows OperationRepository.List results.
type OperationFilter struct {
	AssetID *uint
	From    *time.Time
	To      *time.Time
}

type OperationRepository interface {
	Create(ctx context.Context, op *domain.Operation) error
	GetByID(ctx context.Context, id uint) (*domain.Operation, error)
	List(ctx context.Context, filter OperationFilter) ([]domain.Operation, error)
	Update(ctx context.Context, op *domain.Operation) error
	Delete(ctx context.Context, id uint) error
}
