## syntax=docker/dockerfile:1

FROM golang:1.26-bookworm AS build
WORKDIR /src

# CGO_ENABLED=1 is required: the SQLite driver (mattn/go-sqlite3) wraps the
# C sqlite3 library, so build-essential must be present in this stage.
RUN apt-get update && apt-get install -y --no-install-recommends build-essential \
    && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build -trimpath -ldflags="-s -w" -o /out/api ./cmd/api

FROM debian:bookworm-slim
WORKDIR /app

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates \
    && rm -rf /var/lib/apt/lists/* \
    && useradd --system --create-home --shell /usr/sbin/nologin appuser \
    && mkdir -p /app/data && chown -R appuser:appuser /app

COPY --from=build /out/api /app/api

USER appuser
ENV HTTP_PORT=8080
EXPOSE 8080
VOLUME ["/app/data"]

ENTRYPOINT ["/app/api"]
