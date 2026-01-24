#!/bin/bash

# Memo Studio 一键启动脚本

set -e  # 遇到错误立即退出

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 启动 Memo Studio...${NC}"

# 检查 Go 是否安装
if ! command -v go &> /dev/null; then
    echo -e "${RED}❌ 错误: 未找到 Go，请先安装 Go 1.21+${NC}"
    echo -e "${YELLOW}   安装地址: https://go.dev/dl/${NC}"
    exit 1
fi

# 检查 Node.js 是否安装
if ! command -v node &> /dev/null; then
    echo -e "${RED}❌ 错误: 未找到 Node.js，请先安装 Node.js${NC}"
    echo -e "${YELLOW}   安装地址: https://nodejs.org/${NC}"
    exit 1
fi

# 检查 npm 是否安装
if ! command -v npm &> /dev/null; then
    echo -e "${RED}❌ 错误: 未找到 npm，请先安装 npm${NC}"
    exit 1
fi

# 获取脚本所在目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

# 清理函数
cleanup() {
    echo -e "\n${YELLOW}🛑 正在停止服务...${NC}"
    if [ ! -z "$BACKEND_PID" ]; then
        kill $BACKEND_PID 2>/dev/null || true
    fi
    if [ ! -z "$FRONTEND_PID" ]; then
        kill $FRONTEND_PID 2>/dev/null || true
    fi
    # 清理所有相关进程
    pkill -f "go run main.go" 2>/dev/null || true
    pkill -f "vite dev" 2>/dev/null || true
    echo -e "${GREEN}✅ 服务已停止${NC}"
    exit 0
}

trap cleanup INT TERM

# 检查端口是否被占用
check_port() {
    local port=$1
    if lsof -Pi :$port -sTCP:LISTEN -t >/dev/null 2>&1 ; then
        echo -e "${YELLOW}⚠️  端口 $port 已被占用，尝试停止占用该端口的进程...${NC}"
        lsof -ti:$port | xargs kill -9 2>/dev/null || true
        sleep 1
    fi
}

# 检查并清理端口
check_port 9000
check_port 9001

# 启动后端
echo -e "${BLUE}📦 启动后端服务...${NC}"
cd backend

# 安装 Go 依赖
if [ ! -f "go.sum" ]; then
    echo -e "${YELLOW}📥 安装 Go 依赖...${NC}"
    go mod download
fi

# 检查数据库目录
if [ ! -f "notes.db" ]; then
    echo -e "${YELLOW}💾 数据库文件不存在，将在首次运行时自动创建...${NC}"
fi

# 启动后端（后台运行）
go run main.go > ../backend.log 2>&1 &
BACKEND_PID=$!
cd ..

# 等待后端启动
echo -e "${YELLOW}⏳ 等待后端服务启动...${NC}"
for i in {1..30}; do
    if curl -s http://localhost:9000/api/auth/login > /dev/null 2>&1 || curl -s http://localhost:9000/api/notes > /dev/null 2>&1; then
        echo -e "${GREEN}✅ 后端服务已启动${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}❌ 后端服务启动超时，请检查 backend.log${NC}"
        kill $BACKEND_PID 2>/dev/null || true
        exit 1
    fi
    sleep 1
done

# 启动前端
echo -e "${BLUE}🎨 启动前端应用...${NC}"
cd frontend

# 安装前端依赖
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}📥 安装前端依赖（这可能需要几分钟）...${NC}"
    npm install
fi

# 启动前端（后台运行）
npm run dev > ../frontend.log 2>&1 &
FRONTEND_PID=$!
cd ..

# 等待前端启动
echo -e "${YELLOW}⏳ 等待前端应用启动...${NC}"
for i in {1..30}; do
    if curl -s http://localhost:9001 > /dev/null 2>&1; then
        echo -e "${GREEN}✅ 前端应用已启动${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}❌ 前端应用启动超时，请检查 frontend.log${NC}"
        kill $FRONTEND_PID 2>/dev/null || true
        exit 1
    fi
    sleep 1
done

# 显示启动信息
echo ""
echo -e "${GREEN}════════════════════════════════════════${NC}"
echo -e "${GREEN}✅ Memo Studio 启动成功！${NC}"
echo -e "${GREEN}════════════════════════════════════════${NC}"
echo -e "${BLUE}📝 后端服务: ${GREEN}http://localhost:9000${NC}"
echo -e "${BLUE}🌐 前端应用: ${GREEN}http://localhost:9001${NC}"
echo ""
echo -e "${YELLOW}💡 提示:${NC}"
echo -e "   - 首次使用请先注册账号"
echo -e "   - 日志文件: backend.log 和 frontend.log"
echo -e "   - 按 ${RED}Ctrl+C${NC} 停止所有服务"
echo ""

# 等待用户中断
wait
