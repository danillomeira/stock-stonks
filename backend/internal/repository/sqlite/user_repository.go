package sqlite

import (
	"context"
	"errors"

	"github.com/danillomeira/stock-stonks/backend/internal/config"
	"github.com/danillomeira/stock-stonks/backend/internal/database"
	"github.com/danillomeira/stock-stonks/backend/internal/domain"
	"github.com/danillomeira/stock-stonks/backend/internal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	db    *gorm.DB
	retry config.RetryConfig
}

func NewUserRepository(db *gorm.DB, retry config.RetryConfig) *UserRepository {
	return &UserRepository{db: db, retry: retry}
}

// UpsertByGoogleID inserts the user on first login, or refreshes the cached
// profile fields on every later login. The ID (and therefore every Asset/
// Operation already owned by this user) never changes across logins.
func (r *UserRepository) UpsertByGoogleID(ctx context.Context, user *domain.User) error {
	return database.WithRetry(ctx, r.retry, func() error {
		return r.db.WithContext(ctx).
			Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "google_id"}},
				DoUpdates: clause.AssignmentColumns([]string{"email", "name", "avatar_url", "updated_at"}),
			}).
			Create(user).Error
	})
}

func (r *UserRepository) GetByID(ctx context.Context, id uint) (*domain.User, error) {
	var user domain.User
	err := database.WithRetry(ctx, r.retry, func() error {
		return r.db.WithContext(ctx).First(&user, id).Error
	})
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, repository.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &user, nil
}
