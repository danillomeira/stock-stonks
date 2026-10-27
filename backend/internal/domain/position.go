package domain

import "github.com/shopspring/decimal"

// Position is the running holding state for one asset: how many units are
// held, at what weighted average cost, and profit/loss already realized by
// past sells. It is never persisted - it is always recomputed from
// Operations - because the "custo médio" (weighted average cost) method
// required by Brazilian IR rules must replay history in date order to be
// correct, and a stored snapshot would drift the moment a past Operation is
// edited or deleted.
type Position struct {
	AssetID        uint
	Quantity       decimal.Decimal
	AverageCost    decimal.Decimal
	RealizedProfit decimal.Decimal
}

// ApplyOperation folds one Operation (which must belong to AssetID and be
// processed in ascending date order) into the position using the weighted
// average cost method: buys raise the average cost, sells realize
// profit/loss against the current average cost and never change it.
func (p *Position) ApplyOperation(op Operation) {
	switch op.Type {
	case OperationBuy:
		totalCostBefore := p.AverageCost.Mul(p.Quantity)
		totalCostNow := totalCostBefore.Add(op.NetTotal())
		p.Quantity = p.Quantity.Add(op.Quantity)
		if p.Quantity.IsPositive() {
			p.AverageCost = totalCostNow.Div(p.Quantity)
		}
	case OperationSell:
		costOfSold := p.AverageCost.Mul(op.Quantity)
		p.RealizedProfit = p.RealizedProfit.Add(op.NetTotal().Sub(costOfSold))
		p.Quantity = p.Quantity.Sub(op.Quantity)
	}
}
