package sqlite

import (
	"fmt"

	"github.com/danillomeira/stock-stonks/backend/internal/domain"
	"gorm.io/gorm"
)

// AutoMigrate creates/updates tables for all domain entities. GORM's
// AutoMigrate is adequate for the SQLite MVP; a move to golang-migrate with
// versioned SQL files is the planned path once a networked database (see
// database.Open) is introduced, so schema changes stop being implicit.
func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&domain.User{}, &domain.Asset{}, &domain.Operation{}); err != nil {
		return fmt.Errorf("sqlite: automigrate: %w", err)
	}
	return nil
}
