# 0004. Use shopspring/decimal for all monetary and quantity fields

## Status
Accepted

## Context
This project's purpose is tax reporting. Brazilian IR calculations must
reproduce exact centavo-level amounts across potentially thousands of trades
over many years. `float64` arithmetic accumulates binary rounding error
(e.g. `0.1 + 0.2 != 0.3`), which is unacceptable when the output feeds a tax
declaration.

## Decision
Every monetary or share-quantity field on `domain.Operation`
(`Quantity`, `UnitPrice`, `Fees`) and any value derived from them
(`GrossTotal`, `NetTotal`, `domain.Position.AverageCost`/`RealizedProfit`)
uses `github.com/shopspring/decimal`, never `float64`. Values are
transported as decimal strings over JSON (e.g. `"30.05"`), not floats.

## Consequences
- Arithmetic (average cost, realized profit) is exact, matching what a
  human doing the math by hand would get.
- Adds a small dependency and requires API consumers (the frontend) to
  treat these fields as strings, not numbers, to avoid re-introducing
  float imprecision on the client side.
- GORM column type is declared as `decimal(20,8)` to store the full
  precision rather than relying on SQLite's dynamic typing.

## Alternatives considered
- **`float64`** — rejected: unacceptable rounding risk for tax math.
- **Integer cents (`int64`)** — viable and slightly simpler, but doesn't
  cleanly handle fractional share quantities (B3's mercado fracionário) the
  same way; `decimal.Decimal` handles both money and fractional quantities
  with one type.
