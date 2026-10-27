## syntax=docker/dockerfile:1

FROM node:24-bookworm-slim AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci

FROM node:24-bookworm-slim AS build
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .
RUN npm run build

FROM node:24-bookworm-slim AS runtime
WORKDIR /app
ENV NODE_ENV=production

RUN useradd --system --create-home --shell /usr/sbin/nologin appuser

# next.config.ts has output: "standalone" - this copies only the files
# needed to run, not the full node_modules tree.
COPY --from=build /app/public ./public
COPY --from=build /app/.next/standalone ./
COPY --from=build /app/.next/static ./.next/static

USER appuser
EXPOSE 3000
ENV PORT=3000

CMD ["node", "server.js"]
