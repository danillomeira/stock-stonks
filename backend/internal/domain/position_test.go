package domain_test

import (
	"testing"
	"time"

	"github.com/danillomeira/stock-stonks/backend/internal/domain"
	"github.com/shopspring/decimal"
)

func dec(s string) decimal.Decimal {
	d, err := decimal.NewFromString(s)
	if err != nil {
		panic(err)
	}
	return d
}

func TestPosition_ApplyOperation_WeightedAverageCost(t *testing.T) {
	pos := domain.Position{AssetID: 1}

	// Buy 100 @ 30.00 + 5.00 fees -> avg cost = 3005/100 = 30.05
	pos.ApplyOperation(domain.Operation{
		AssetID: 1, Type: domain.OperationBuy, Date: time.Now(),
		Quantity: dec("100"), UnitPrice: dec("30.00"), Fees: dec("5.00"),
	})
	if !pos.Quantity.Equal(dec("100")) {
		t.Fatalf("quantity = %s, want 100", pos.Quantity)
	}
	if !pos.AverageCost.Equal(dec("30.05")) {
		t.Fatalf("average cost = %s, want 30.05", pos.AverageCost)
	}

	// Buy another 100 @ 31.00, no fees -> total cost 3005+3100=6105 / 200 = 30.525
	pos.ApplyOperation(domain.Operation{
		AssetID: 1, Type: domain.OperationBuy, Date: time.Now(),
		Quantity: dec("100"), UnitPrice: dec("31.00"), Fees: dec("0"),
	})
	if !pos.AverageCost.Equal(dec("30.525")) {
		t.Fatalf("average cost after 2nd buy = %s, want 30.525", pos.AverageCost)
	}

	// Sell 40 @ 35.00 - 3.00 fees -> net = 1400-3=1397; cost of sold = 30.525*40=1221
	// realized profit = 1397-1221 = 176; average cost must NOT change on sell.
	pos.ApplyOperation(domain.Operation{
		AssetID: 1, Type: domain.OperationSell, Date: time.Now(),
		Quantity: dec("40"), UnitPrice: dec("35.00"), Fees: dec("3.00"),
	})
	if !pos.Quantity.Equal(dec("160")) {
		t.Fatalf("quantity after sell = %s, want 160", pos.Quantity)
	}
	if !pos.AverageCost.Equal(dec("30.525")) {
		t.Fatalf("average cost must be unchanged by a sell, got %s", pos.AverageCost)
	}
	if !pos.RealizedProfit.Equal(dec("176")) {
		t.Fatalf("realized profit = %s, want 176", pos.RealizedProfit)
	}
}

func TestOperation_Validate(t *testing.T) {
	cases := []struct {
		name string
		op   domain.Operation
		want error
	}{
		{"missing user", domain.Operation{AssetID: 1, Type: domain.OperationBuy, Quantity: dec("1"), UnitPrice: dec("1"), Date: time.Now()}, domain.ErrInvalidUserID},
		{"missing asset", domain.Operation{UserID: 1, Type: domain.OperationBuy, Quantity: dec("1"), UnitPrice: dec("1"), Date: time.Now()}, domain.ErrInvalidAssetID},
		{"bad type", domain.Operation{UserID: 1, AssetID: 1, Type: "HOLD", Quantity: dec("1"), UnitPrice: dec("1"), Date: time.Now()}, domain.ErrInvalidOperation},
		{"zero quantity", domain.Operation{UserID: 1, AssetID: 1, Type: domain.OperationBuy, Quantity: dec("0"), UnitPrice: dec("1"), Date: time.Now()}, domain.ErrInvalidQuantity},
		{"zero price", domain.Operation{UserID: 1, AssetID: 1, Type: domain.OperationBuy, Quantity: dec("1"), UnitPrice: dec("0"), Date: time.Now()}, domain.ErrInvalidUnitPrice},
		{"negative fees", domain.Operation{UserID: 1, AssetID: 1, Type: domain.OperationBuy, Quantity: dec("1"), UnitPrice: dec("1"), Fees: dec("-1"), Date: time.Now()}, domain.ErrNegativeFees},
		{"zero date", domain.Operation{UserID: 1, AssetID: 1, Type: domain.OperationBuy, Quantity: dec("1"), UnitPrice: dec("1")}, domain.ErrInvalidDate},
		{"valid", domain.Operation{UserID: 1, AssetID: 1, Type: domain.OperationBuy, Quantity: dec("1"), UnitPrice: dec("1"), Date: time.Now()}, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.op.Validate(); err != tc.want {
				t.Fatalf("Validate() = %v, want %v", err, tc.want)
			}
		})
	}
}
