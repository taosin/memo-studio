### Build SvelteKit static
FROM node:20-bookworm AS kit-builder
WORKDIR /app
# 更新系统包以修复安全漏洞
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
COPY kit/package.json kit/package-lock.json ./kit/
RUN cd kit && npm ci
COPY kit ./kit
RUN cd kit && npm run build

### Build Go binary (CGO + sqlite_fts5)
FROM golang:1.23-bookworm AS go-builder
WORKDIR /app
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends build-essential ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*
# 设置 Go 代理（可选，有助于解决网络问题）
ENV GOPROXY=https://proxy.golang.org,direct
COPY backend/go.mod backend/go.sum ./backend/
RUN cd backend && go mod download
COPY backend ./backend
COPY --from=kit-builder /app/kit/build ./backend/public
RUN cd backend && CGO_ENABLED=1 go build -tags sqlite_fts5 -o /out/memo-studio .

### Runtime
FROM debian:bookworm-slim
WORKDIR /app
# 更新所有包到最新版本以修复安全漏洞
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends ca-certificates tzdata wget && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# non-root user
RUN useradd -m -u 10001 appuser

# 创建数据目录并设置权限（在切换用户之前）
RUN mkdir -p /data && chown -R appuser:appuser /data

COPY --from=go-builder /out/memo-studio /app/memo-studio

# 切换到非 root 用户
USER appuser

ENV GIN_MODE=release \
    PORT=9000 \
    MEMO_DB_PATH=/data/notes.db \
    MEMO_STORAGE_DIR=/data/storage

EXPOSE 9000
VOLUME ["/data"]

HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
    CMD wget -qO- "http://127.0.0.1:${PORT:-9000}/health" >/dev/null 2>&1 || exit 1

CMD ["/app/memo-studio"]

