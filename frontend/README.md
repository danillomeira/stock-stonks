# Stock-Stonks — frontend

Next.js (App Router, TypeScript) app for the Stock-Stonks project. See the
[root README](../README.md) for the full project overview and
[`agents.md`](../agents.md) for AI-agent working conventions.

## Development

```sh
cp .env.example .env.local   # set NEXT_PUBLIC_API_URL if the backend isn't on :8080
npm install
npm run dev
```

## Pages

- `/assets` — register tickers (STOCK/FII/ETF/BDR).
- `/operations` — register buy/sell trades.
- `/positions` — current quantity/average cost/realized profit per asset,
  computed by the backend from the full operation history.

## Checks

```sh
npm run lint
npx tsc --noEmit
npm run build
```
