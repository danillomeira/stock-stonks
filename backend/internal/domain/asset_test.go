package domain_test

import (
	"testing"

	"github.com/danillomeira/stock-stonks/backend/internal/domain"
)

func TestAsset_Validate(t *testing.T) {
	cases := []struct {
		name  string
		asset domain.Asset
		want  error
	}{
		{"missing user", domain.Asset{Ticker: "PETR4", Name: "Petrobras", Type: domain.AssetTypeStock}, domain.ErrInvalidUserID},
		{"empty ticker", domain.Asset{UserID: 1, Name: "Petrobras", Type: domain.AssetTypeStock}, domain.ErrInvalidTicker},
		{"empty name", domain.Asset{UserID: 1, Ticker: "PETR4", Type: domain.AssetTypeStock}, domain.ErrInvalidName},
		{"invalid type", domain.Asset{UserID: 1, Ticker: "PETR4", Name: "Petrobras", Type: "CRYPTO"}, domain.ErrInvalidType},
		{"valid", domain.Asset{UserID: 1, Ticker: "PETR4", Name: "Petrobras", Type: domain.AssetTypeStock}, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.asset.Validate(); err != tc.want {
				t.Fatalf("Validate() = %v, want %v", err, tc.want)
			}
		})
	}
}
