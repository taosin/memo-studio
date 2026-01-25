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

./dist/memo-studio

