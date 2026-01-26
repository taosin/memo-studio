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

open_url() {
    local url="$1"
    if command -v open &> /dev/null; then
        open "$url" >/dev/null 2>&1 || true
    elif command -v xdg-open &> /dev/null; then
        xdg-open "$url" >/dev/null 2>&1 || true
    elif command -v cmd.exe &> /dev/null; then
        cmd.exe /c start "$url" >/dev/null 2>&1 || true
    else
        echo -e "${YELLOW}请手动打开: ${url}${NC}"
    fi
}

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
    pkill -f "go run .*sqlite_fts5" 2>/dev/null || true
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
echo -e "${BLUE}🔍 检查端口占用...${NC}"
if ! check_port 9000; then
    echo -e "${RED}❌ 无法清理端口 9000，请手动处理${NC}"
    exit 1
fi
if ! check_port 9001; then
    echo -e "${RED}❌ 无法清理端口 9001，请手动处理${NC}"
    exit 1
fi

# 启动后端
echo -e "${BLUE}📦 启动后端服务...${NC}"
cd backend

# 设置 Go 代理（使用国内镜像，避免网络问题）
if [ -z "$GOPROXY" ]; then
    export GOPROXY=https://goproxy.cn,direct
    echo -e "${YELLOW}📡 已设置 Go 代理为国内镜像: goproxy.cn${NC}"
fi

# 安装 Go 依赖
if [ ! -f "go.sum" ]; then
    echo -e "${YELLOW}📥 安装 Go 依赖...${NC}"
    go mod download && go mod tidy || {
        echo -e "${YELLOW}⚠️  使用备用代理重试...${NC}"
        export GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
        go mod download && go mod tidy || {
            echo -e "${RED}❌ Go 依赖安装失败${NC}"
            echo -e "${YELLOW}   请检查网络连接或手动设置代理:${NC}"
            echo -e "${YELLOW}   export GOPROXY=https://goproxy.cn,direct${NC}"
            echo -e "${YELLOW}   cd backend && go mod download && go mod tidy${NC}"
            exit 1
        }
    }
    echo -e "${GREEN}✅ Go 依赖安装完成${NC}"
fi

# 检查数据库目录
if [ ! -f "notes.db" ]; then
    echo -e "${YELLOW}💾 数据库文件不存在，将在首次运行时自动创建...${NC}"
fi

# 启动后端（后台运行，输出到日志）
echo -e "${YELLOW}⏳ 正在启动后端服务...${NC}"
# 确保在 backend 目录中运行
(cd "$(pwd)" && go run -tags sqlite_fts5 main.go > ../backend.log 2>&1) &
BACKEND_PID=$!
cd ..
# 给后端一点时间开始启动
sleep 2

# 等待后端启动
echo -e "${YELLOW}⏳ 等待后端服务启动（最多等待30秒）...${NC}"
BACKEND_READY=false
for i in {1..30}; do
    # 使用健康检查端点
    if curl -s -f http://localhost:9000/health > /dev/null 2>&1; then
        echo -e "${GREEN}✅ 后端服务已启动${NC}"
        BACKEND_READY=true
        break
    fi
    # 显示进度
    if [ $((i % 5)) -eq 0 ]; then
        echo -e "${YELLOW}   已等待 ${i} 秒...${NC}"
    fi
    sleep 1
done

if [ "$BACKEND_READY" = false ]; then
    echo -e "${RED}❌ 后端服务启动超时${NC}"
    echo -e "${YELLOW}📋 检查后端日志:${NC}"
    if [ -f "backend.log" ]; then
        tail -30 backend.log
    else
        echo "  日志文件不存在，可能后端进程未启动"
    fi
    echo ""
    echo -e "${YELLOW}💡 排查建议:${NC}"
    echo -e "   1. 检查端口 9000 是否被占用: ${BLUE}lsof -i :9000${NC}"
    echo -e "   2. 手动启动后端查看错误: ${BLUE}cd backend && go run -tags sqlite_fts5 main.go${NC}"
    echo -e "   3. 检查数据库文件权限: ${BLUE}ls -la backend/notes.db${NC}"
    kill $BACKEND_PID 2>/dev/null || true
    exit 1
fi

# 检查后端进程是否还在运行
if ! kill -0 $BACKEND_PID 2>/dev/null; then
    echo -e "${RED}❌ 后端服务启动失败（进程已退出）${NC}"
    echo -e "${YELLOW}📋 后端日志:${NC}"
    cat backend.log 2>/dev/null || echo "无法读取日志文件"
    exit 1
fi

# 启动前端
echo -e "${BLUE}🎨 启动前端应用...${NC}"
cd frontend

# 安装前端依赖
if [ ! -d "node_modules" ]; then
    echo -e "${YELLOW}📥 安装前端依赖（这可能需要几分钟）...${NC}"
    npm install || {
        echo -e "${RED}❌ 前端依赖安装失败${NC}"
        exit 1
    }
fi

# 启动前端（后台运行，输出到日志）
echo -e "${YELLOW}⏳ 正在启动前端应用...${NC}"
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
        echo -e "${RED}❌ 前端应用启动超时${NC}"
        echo -e "${YELLOW}📋 前端日志（最后20行）:${NC}"
        tail -20 frontend.log 2>/dev/null || echo "无法读取日志文件"
        kill $FRONTEND_PID 2>/dev/null || true
        exit 1
    fi
    sleep 1
done

# 检查前端是否有错误
if ! kill -0 $FRONTEND_PID 2>/dev/null; then
    echo -e "${RED}❌ 前端应用启动失败${NC}"
    echo -e "${YELLOW}📋 前端日志:${NC}"
    cat frontend.log 2>/dev/null || echo "无法读取日志文件"
    exit 1
fi

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
echo -e "   - 查看日志: tail -f backend.log 或 tail -f frontend.log"
echo -e "   - 按 ${RED}Ctrl+C${NC} 停止所有服务"
echo ""

# 自动打开浏览器（旧前端无路由跳转，这里打开首页即可）
open_url "http://localhost:9001"

# 等待用户中断
wait
