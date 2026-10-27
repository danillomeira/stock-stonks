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

// Every read/update/delete below takes userID explicitly and scopes its
// query by it (WHERE user_id = ? AND ...), rather than trusting the ID
// alone - a row that exists but belongs to a different user is
// indistinguishable from one that doesn't exist (ErrNotFound either way),
// so ownership can never leak through a different error shape.

type UserRepository interface {
	// UpsertByGoogleID creates the user on first login, or updates the
	// cached profile fields (email/name/avatar) on every later login.
	UpsertByGoogleID(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uint) (*domain.User, error)
}

type AssetRepository interface {
	Create(ctx context.Context, asset *domain.Asset) error
	GetByID(ctx context.Context, userID, id uint) (*domain.Asset, error)
	GetByTicker(ctx context.Context, userID uint, ticker string) (*domain.Asset, error)
	List(ctx context.Context, userID uint) ([]domain.Asset, error)
	Update(ctx context.Context, asset *domain.Asset) error
	Delete(ctx context.Context, userID, id uint) error
}

// OperationFilter narrows OperationRepository.List results, always within
// the caller's own data (see UserID).
type OperationFilter struct {
	UserID  uint
	AssetID *uint
	From    *time.Time
	To      *time.Time
}

type OperationRepository interface {
	Create(ctx context.Context, op *domain.Operation) error
	GetByID(ctx context.Context, userID, id uint) (*domain.Operation, error)
	List(ctx context.Context, filter OperationFilter) ([]domain.Operation, error)
	Update(ctx context.Context, op *domain.Operation) error
	Delete(ctx context.Context, userID, id uint) error
}
