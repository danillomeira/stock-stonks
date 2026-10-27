# Stock-Stonks

Acompanhamento e registro de operações em ações (e FIIs/ETFs/BDRs), com
exportação de dados para o Imposto de Renda, construído em etapas a partir
de um MVP.

## Status

MVP: login com Google (multiusuário — cada conta vê só os próprios dados),
cadastro de ativos, registro de operações (compra/venda) e cálculo de
posição (preço médio, lucro realizado) pelo método de custo médio ponderado.
Cotação em tempo real, exportação para IR e classificação de day-trade ainda
não implementadas — veja `specs/` para o roadmap.

## Stack

- **Backend**: Go + GORM + SQLite (`backend/`) — veja
  `docs/adr/0003-gorm-sqlite-with-retry-backoff-seam.md` para a arquitetura
  de retry/backoff e o caminho para bancos mais escaláveis.
- **Frontend**: Next.js (App Router, TypeScript) (`frontend/`).
- **CI/CD**: GitHub Actions — lint/test/build em cada push, build de imagens
  Docker publicadas no GHCR a partir da `main` (sem deploy automático ainda,
  veja `docs/adr/0007-ci-only-pipeline-no-deploy-target-yet.md`).

## Rodando localmente

### Backend

```sh
cd backend
go run ./cmd/api
# http://localhost:8080 — GET /health
```

Variáveis de ambiente (todas com default razoável para dev, ver
`backend/internal/config/config.go`): `HTTP_PORT`, `DB_DRIVER`, `DB_DSN`,
`CORS_ALLOW_ORIGIN`, `DB_RETRY_MAX_ATTEMPTS`, `DB_RETRY_BASE_DELAY`,
`DB_RETRY_MAX_DELAY`. Sem `GOOGLE_CLIENT_ID`/`GOOGLE_CLIENT_SECRET` o
servidor roda normalmente — só o login fica indisponível (ver abaixo).

#### Login com Google

1. No [Google Cloud Console](https://console.cloud.google.com/apis/credentials),
   crie um **OAuth client ID** do tipo "Web application".
2. Em "Authorized redirect URIs", adicione
   `http://localhost:8080/api/v1/auth/google/callback`.
3. Exporte as credenciais antes de rodar o backend:
   ```sh
   export GOOGLE_CLIENT_ID=...
   export GOOGLE_CLIENT_SECRET=...
   ```
4. (Opcional) `SESSION_SECRET=<valor fixo>` para a sessão sobreviver a um
   restart do servidor — sem isso, um segredo aleatório é gerado a cada
   execução e todo mundo precisa logar de novo (ver ADR 0009).

### Frontend

```sh
cd frontend
cp .env.example .env.local   # ajuste NEXT_PUBLIC_API_URL se necessário
npm install
npm run dev
# http://localhost:3000
```

### Com containers (Docker ou Podman)

O `docker-compose.yml` é agnóstico de runtime — funciona com `docker
compose` ou `podman-compose`/`podman compose`:

```sh
podman-compose up --build
# ou: docker compose up --build
```

## Testes

```sh
cd backend && go test ./... -race -cover
cd frontend && npm run lint && npx tsc --noEmit && npm run build
```

## Estrutura do repositório

```
backend/            API Go (domain → repository → service → httpserver → cmd)
frontend/           App Next.js
specs/              Specs orientadas a feature (spec-kit style) + constitution.md
docs/adr/           Architecture Decision Records
.claude/skills/     Skills para agentes de IA (find-skills, adr, documentation, reuse)
docs/catalog.md     Índice de componentes/helpers reutilizáveis (ver skill `reuse`)
.github/workflows/  CI (lint/test/build) + publicação de imagens no GHCR
```

Para trabalhar neste repositório com um agente de IA, comece por
[`agents.md`](agents.md).
