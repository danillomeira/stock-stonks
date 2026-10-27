package domain

import (
	"errors"
	"strings"
	"time"
)

// AssetType classifies an asset for tax-treatment purposes, since Brazilian
// income tax rules differ between common stocks, FIIs and other classes.
type AssetType string

const (
	AssetTypeStock AssetType = "STOCK" // Ação
	AssetTypeFII   AssetType = "FII"   // Fundo de Investimento Imobiliário
	AssetTypeETF   AssetType = "ETF"
	AssetTypeBDR   AssetType = "BDR"
)

func (t AssetType) Valid() bool {
	switch t {
	case AssetTypeStock, AssetTypeFII, AssetTypeETF, AssetTypeBDR:
		return true
	default:
		return false
	}
}

// Asset represents a tradable ticker (e.g. PETR4, HGLG11) owned by one User.
// A ticker may be registered independently by multiple users - there is no
// global uniqueness, only per-user (see the composite index below).
type Asset struct {
	ID        uint      `gorm:"primaryKey"`
	UserID    uint      `gorm:"not null;index;uniqueIndex:idx_assets_user_ticker"`
	Ticker    string    `gorm:"size:16;not null;uniqueIndex:idx_assets_user_ticker"`
	Name      string    `gorm:"size:255;not null"`
	Type      AssetType `gorm:"size:16;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrInvalidUserID = errors.New("domain: user id must be set")
	ErrInvalidTicker = errors.New("domain: ticker must not be empty")
	ErrInvalidName   = errors.New("domain: name must not be empty")
	ErrInvalidType   = errors.New("domain: asset type is invalid")
)

// Validate enforces the invariants an Asset must hold before it is persisted.
func (a *Asset) Validate() error {
	if a.UserID == 0 {
		return ErrInvalidUserID
	}
	if strings.TrimSpace(a.Ticker) == "" {
		return ErrInvalidTicker
	}
	if strings.TrimSpace(a.Name) == "" {
		return ErrInvalidName
	}
	if !a.Type.Valid() {
		return ErrInvalidType
	}
	return nil
}
