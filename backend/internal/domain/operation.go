package domain

import (
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

// OperationType is the trade side. Day-trade classification (which changes
// the applicable Brazilian IR rate) is derived later from same-day buy/sell
// pairs on the same asset, not stored here.
type OperationType string

const (
	OperationBuy  OperationType = "BUY"  // Compra
	OperationSell OperationType = "SELL" // Venda
)

func (t OperationType) Valid() bool {
	return t == OperationBuy || t == OperationSell
}

// Operation is a single buy or sell trade of an Asset on a given date.
// Monetary fields use decimal.Decimal, never float64: IR tax calculations
// must reproduce exact centavo amounts, and binary floats accumulate
// rounding error across thousands of trades.
type Operation struct {
	ID      uint  `gorm:"primaryKey"`
	UserID  uint  `gorm:"not null;index"`
	AssetID uint  `gorm:"not null;index"`
	Asset   Asset `gorm:"foreignKey:AssetID"`

	Date      time.Time       `gorm:"not null;index"`
	Type      OperationType   `gorm:"size:8;not null"`
	Quantity  decimal.Decimal `gorm:"type:decimal(20,8);not null"` // fractional lots are allowed on B3
	UnitPrice decimal.Decimal `gorm:"type:decimal(20,8);not null"`
	Fees      decimal.Decimal `gorm:"type:decimal(20,8);not null"` // corretagem + emolumentos + ISS combined
	Broker    string          `gorm:"size:120"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrInvalidAssetID   = errors.New("domain: asset id must be set")
	ErrInvalidOperation = errors.New("domain: operation type is invalid")
	ErrInvalidQuantity  = errors.New("domain: quantity must be greater than zero")
	ErrInvalidUnitPrice = errors.New("domain: unit price must be greater than zero")
	ErrNegativeFees     = errors.New("domain: fees must not be negative")
	ErrInvalidDate      = errors.New("domain: date must not be zero")
)

// Validate enforces the invariants an Operation must hold before it is persisted.
func (o *Operation) Validate() error {
	if o.UserID == 0 {
		return ErrInvalidUserID
	}
	if o.AssetID == 0 {
		return ErrInvalidAssetID
	}
	if !o.Type.Valid() {
		return ErrInvalidOperation
	}
	if o.Quantity.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidQuantity
	}
	if o.UnitPrice.LessThanOrEqual(decimal.Zero) {
		return ErrInvalidUnitPrice
	}
	if o.Fees.IsNegative() {
		return ErrNegativeFees
	}
	if o.Date.IsZero() {
		return ErrInvalidDate
	}
	return nil
}

// GrossTotal is quantity * unit price, before fees.
func (o *Operation) GrossTotal() decimal.Decimal {
	return o.Quantity.Mul(o.UnitPrice)
}

// NetTotal is the amount actually paid (buy) or received (sell) after fees.
func (o *Operation) NetTotal() decimal.Decimal {
	if o.Type == OperationBuy {
		return o.GrossTotal().Add(o.Fees)
	}
	return o.GrossTotal().Sub(o.Fees)
}
