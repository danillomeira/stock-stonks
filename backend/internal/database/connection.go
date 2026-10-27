package database

import (
	"fmt"

	"github.com/danillomeira/stock-stonks/backend/internal/config"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Open connects to the database selected by cfg.DBDriver and returns a ready
// *gorm.DB. It is the single seam for adding another backend later: a new
// driver only needs a new case here plus its own entry in IsRetryable -
// repositories, services and handlers never see the driver name.
func Open(cfg config.Config) (*gorm.DB, error) {
	switch cfg.DBDriver {
	case "sqlite", "":
		return openSQLite(cfg.DBDSN)
	case "postgres", "mysql":
		return nil, fmt.Errorf("database: driver %q is reserved for future use but not implemented yet; use \"sqlite\" for now", cfg.DBDriver)
	default:
		return nil, fmt.Errorf("database: unknown driver %q", cfg.DBDriver)
	}
}

func openSQLite(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger:         logger.Default.LogMode(logger.Warn),
		TranslateError: true, // lets repositories match gorm.ErrDuplicatedKey instead of parsing driver-specific messages
	})
	if err != nil {
		return nil, fmt.Errorf("database: open sqlite: %w", err)
	}

	// SQLite only allows one writer at a time; a single pooled connection
	// plus the retry/backoff wrapper in retry.go is how concurrent writers
	// are made to behave instead of failing with "database is locked".
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("database: get underlying sql.DB: %w", err)
	}
	sqlDB.SetMaxOpenConns(1)

	return db, nil
}
