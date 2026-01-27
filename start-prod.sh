#!/bin/bash

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

if [ ! -f "dist/memo-studio" ]; then
  echo -e "${YELLOW}未发现 dist/memo-studio，开始构建...${NC}"
  ./build-prod.sh
fi

echo -e "${BLUE}🚀 启动生产模式（Go 托管前端静态）...${NC}"
echo -e "${BLUE}🌐 打开: ${GREEN}http://localhost:9000${NC}"
echo -e "${YELLOW}数据库路径可通过 MEMO_DB_PATH 指定（默认 ./notes.db）${NC}"

open_url() {
  local url="$1"
  if command -v open >/dev/null 2>&1; then
    open "$url" >/dev/null 2>&1 || true
  elif command -v xdg-open >/dev/null 2>&1; then
    xdg-open "$url" >/dev/null 2>&1 || true
  elif command -v cmd.exe >/dev/null 2>&1; then
    cmd.exe /c start "$url" >/dev/null 2>&1 || true
  else
    echo -e "${YELLOW}请手动打开: ${url}${NC}"
  fi
}

cleanup() {
  if [ -n "${PID:-}" ]; then
    kill "$PID" 2>/dev/null || true
  fi
  exit 0
}
trap cleanup INT TERM

./dist/memo-studio &
PID=$!

# 等待后端就绪后再打开（优先打开登录页）
READY=false
for i in {1..60}; do
  if curl -s -f "http://localhost:${PORT:-9000}/health" >/dev/null 2>&1; then
    READY=true
    break
  fi
  sleep 1
done
if [ "$READY" = true ]; then
  open_url "http://localhost:${PORT:-9000}/login"
else
  open_url "http://localhost:${PORT:-9000}/"
fi

wait "$PID"

