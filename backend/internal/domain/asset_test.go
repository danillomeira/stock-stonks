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
		{"empty ticker", domain.Asset{Name: "Petrobras", Type: domain.AssetTypeStock}, domain.ErrInvalidTicker},
		{"empty name", domain.Asset{Ticker: "PETR4", Type: domain.AssetTypeStock}, domain.ErrInvalidName},
		{"invalid type", domain.Asset{Ticker: "PETR4", Name: "Petrobras", Type: "CRYPTO"}, domain.ErrInvalidType},
		{"valid", domain.Asset{Ticker: "PETR4", Name: "Petrobras", Type: domain.AssetTypeStock}, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.asset.Validate(); err != tc.want {
				t.Fatalf("Validate() = %v, want %v", err, tc.want)
			}
		})
	}
}
