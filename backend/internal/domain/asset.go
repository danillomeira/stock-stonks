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

// Asset represents a tradable ticker (e.g. PETR4, HGLG11).
type Asset struct {
	ID        uint      `gorm:"primaryKey"`
	Ticker    string    `gorm:"uniqueIndex;size:16;not null"`
	Name      string    `gorm:"size:255;not null"`
	Type      AssetType `gorm:"size:16;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrInvalidTicker = errors.New("domain: ticker must not be empty")
	ErrInvalidName   = errors.New("domain: name must not be empty")
	ErrInvalidType   = errors.New("domain: asset type is invalid")
)

// Validate enforces the invariants an Asset must hold before it is persisted.
func (a *Asset) Validate() error {
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
