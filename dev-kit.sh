#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

cleanup() {
  echo -e "\n${YELLOW}🛑 正在停止 dev 服务...${NC}"
  if [ -n "${BACKEND_PID:-}" ]; then kill "$BACKEND_PID" 2>/dev/null || true; fi
  if [ -n "${KIT_PID:-}" ]; then kill "$KIT_PID" 2>/dev/null || true; fi
  pkill -f "go run .*sqlite_fts5" 2>/dev/null || true
  pkill -f "vite dev" 2>/dev/null || true
  echo -e "${GREEN}✅ 已停止${NC}"
}
trap cleanup INT TERM

echo -e "${BLUE}🚀 Dev 启动（Go + SvelteKit）...${NC}"

command -v go >/dev/null 2>&1 || { echo -e "${RED}❌ 需要安装 Go${NC}"; exit 1; }
command -v node >/dev/null 2>&1 || { echo -e "${RED}❌ 需要安装 Node.js${NC}"; exit 1; }
command -v npm >/dev/null 2>&1 || { echo -e "${RED}❌ 需要安装 npm${NC}"; exit 1; }

check_port() {
  local port=$1
  if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1 ; then
    echo -e "${YELLOW}⚠️  端口 $port 已被占用，尝试停止占用进程...${NC}"
    lsof -ti:$port | xargs kill -9 2>/dev/null || true
    sleep 1
  fi
}

echo -e "${BLUE}🔍 检查端口占用...${NC}"
check_port 9000
check_port 9001

echo -e "${BLUE}📦 启动后端（:9000，FTS5）...${NC}"
(cd backend && go run -tags sqlite_fts5 . > ../backend.log 2>&1) &
BACKEND_PID=$!

# 等待后端健康检查
BACKEND_READY=false
for i in {1..20}; do
  if curl -s -f http://localhost:9000/health > /dev/null 2>&1; then
    BACKEND_READY=true
    break
  fi
  sleep 1
done
if [ "$BACKEND_READY" = false ]; then
  echo -e "${RED}❌ 后端未成功启动（/health 不可达）${NC}"
  echo -e "${YELLOW}📋 backend.log（最后50行）:${NC}"
  tail -50 backend.log 2>/dev/null || true
  exit 1
fi

echo -e "${BLUE}🎨 启动前端（SvelteKit dev :9001）...${NC}"
if [ ! -d "kit/node_modules" ]; then
  echo -e "${YELLOW}📥 安装 kit 依赖（首次较慢）...${NC}"
  (cd kit && npm install)
fi
(cd kit && npm run dev > ../kit.log 2>&1) &
KIT_PID=$!

# 等待前端首页可达
KIT_READY=false
for i in {1..20}; do
  if curl -s http://localhost:9001/ > /dev/null 2>&1; then
    KIT_READY=true
    break
  fi
  sleep 1
done
if [ "$KIT_READY" = false ]; then
  echo -e "${RED}❌ 前端未成功启动（9001 不可达）${NC}"
  echo -e "${YELLOW}📋 kit.log（最后80行）:${NC}"
  tail -80 kit.log 2>/dev/null || true
  exit 1
fi

echo -e "${GREEN}✅ Dev 已启动${NC}"
echo -e "${BLUE}📝 API: ${GREEN}http://localhost:9000/api${NC}"
echo -e "${BLUE}🌐 Web: ${GREEN}http://localhost:9001${NC}"
echo -e "${YELLOW}查看日志: tail -f backend.log 或 tail -f kit.log${NC}"
echo -e "${YELLOW}按 Ctrl+C 停止${NC}"

wait

