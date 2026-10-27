package sqlite_test

import (
	"context"
	"testing"
	"time"

	"github.com/danillomeira/stock-stonks/backend/internal/config"
	"github.com/danillomeira/stock-stonks/backend/internal/domain"
	"github.com/danillomeira/stock-stonks/backend/internal/repository/sqlite"
	gormsqlite "gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(gormsqlite.Open("file::memory:?cache=shared"), &gorm.Config{TranslateError: true})
	if err != nil {
		t.Fatalf("open in-memory sqlite: %v", err)
	}
	if err := sqlite.AutoMigrate(db); err != nil {
		t.Fatalf("automigrate: %v", err)
	}
	return db
}

func testRetryConfig() config.RetryConfig {
	return config.RetryConfig{MaxAttempts: 1, BaseDelay: time.Millisecond, MaxDelay: time.Millisecond}
}

func TestUserRepository_UpsertByGoogleID_CreatesThenUpdatesSameRow(t *testing.T) {
	db := newTestDB(t)
	repo := sqlite.NewUserRepository(db, testRetryConfig())
	ctx := context.Background()

	first := &domain.User{GoogleID: "google-123", Email: "old@example.com", Name: "Old Name"}
	if err := repo.UpsertByGoogleID(ctx, first); err != nil {
		t.Fatalf("first upsert: %v", err)
	}
	if first.ID == 0 {
		t.Fatalf("expected first upsert to assign an ID, got 0")
	}

	second := &domain.User{GoogleID: "google-123", Email: "new@example.com", Name: "New Name", AvatarURL: "https://example.com/a.png"}
	if err := repo.UpsertByGoogleID(ctx, second); err != nil {
		t.Fatalf("second upsert: %v", err)
	}

	got, err := repo.GetByID(ctx, first.ID)
	if err != nil {
		t.Fatalf("GetByID: %v", err)
	}
	if got.Email != "new@example.com" || got.Name != "New Name" || got.AvatarURL != "https://example.com/a.png" {
		t.Fatalf("expected profile fields refreshed, got %+v", got)
	}

	var count int64
	if err := db.Model(&domain.User{}).Where("google_id = ?", "google-123").Count(&count).Error; err != nil {
		t.Fatalf("count: %v", err)
	}
	if count != 1 {
		t.Fatalf("expected exactly 1 row for google_id, got %d (upsert created a duplicate instead of updating)", count)
	}
}
