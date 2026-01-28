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

# 设置 Go 环境变量（优化下载和缓存）
ENV GOPROXY=https://proxy.golang.org,direct \
    GOMODCACHE=/go/pkg/mod \
    GOCACHE=/go/cache

# 安装构建依赖
RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y --no-install-recommends build-essential ca-certificates && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# 先复制 go.mod 和 go.sum（利用 Docker 层缓存）
# 只有当依赖文件变化时才重新下载
COPY backend/go.mod backend/go.sum ./backend/

# 下载依赖（这一层会被缓存，除非 go.mod/go.sum 变化）
RUN cd backend && \
    go mod download && \
    go mod verify

# 复制源代码（放在依赖下载之后，避免代码变化导致重新下载依赖）
COPY backend ./backend
COPY --from=kit-builder /app/kit/build ./backend/public

# 构建应用
RUN cd backend && \
    CGO_ENABLED=1 go build -tags sqlite_fts5 -o /out/memo-studio .

### Runtime
FROM debian:stable-20260112-slim
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

