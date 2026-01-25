### Build SvelteKit static
FROM node:20-bookworm AS kit-builder
WORKDIR /app
COPY kit/package.json kit/package-lock.json ./kit/
RUN cd kit && npm ci
COPY kit ./kit
RUN cd kit && npm run build

### Build Go binary (CGO + sqlite_fts5)
FROM golang:1.21-bookworm AS go-builder
WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends build-essential ca-certificates && rm -rf /var/lib/apt/lists/*
COPY backend/go.mod backend/go.sum ./backend/
RUN cd backend && go mod download
COPY backend ./backend
COPY --from=kit-builder /app/kit/build ./backend/public
RUN cd backend && CGO_ENABLED=1 go build -tags sqlite_fts5 -o /out/memo-studio .

### Runtime
FROM debian:bookworm-slim
WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates tzdata wget && rm -rf /var/lib/apt/lists/*

# non-root user
RUN useradd -m -u 10001 appuser
USER appuser

COPY --from=go-builder /out/memo-studio /app/memo-studio

ENV GIN_MODE=release \
    PORT=9000 \
    MEMO_DB_PATH=/data/notes.db \
    MEMO_STORAGE_DIR=/data/storage

EXPOSE 9000
VOLUME ["/data"]

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget -qO- "http://127.0.0.1:${PORT:-9000}/health" >/dev/null 2>&1 || exit 1

CMD ["/app/memo-studio"]

