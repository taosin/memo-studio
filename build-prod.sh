#!/bin/bash

set -e

GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

echo -e "${BLUE}🏗️  构建前端（SvelteKit 静态）...${NC}"
command -v npm >/dev/null 2>&1 || { echo "need npm"; exit 1; }
(cd kit && npm ci && npm run build)

echo -e "${BLUE}📦 同步前端产物到 backend/public ...${NC}"
mkdir -p backend/public
if command -v rsync >/dev/null 2>&1; then
  rsync -a --delete kit/build/ backend/public/
else
  rm -rf backend/public/*
  cp -R kit/build/* backend/public/
fi

echo -e "${BLUE}🔧 构建 Go 二进制（启用 FTS5）...${NC}"
mkdir -p dist
(cd backend && go build -tags sqlite_fts5 -o ../dist/memo-studio .)

echo -e "${GREEN}✅ 构建完成${NC}"
echo -e "${YELLOW}运行: ./dist/memo-studio${NC}"

