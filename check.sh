#!/bin/bash

# 诊断脚本 - 检查常见问题

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}🔍 Memo Studio 诊断工具${NC}"
echo ""

# 检查 Go
echo -e "${BLUE}检查 Go 环境...${NC}"
if command -v go &> /dev/null; then
    GO_VERSION=$(go version | awk '{print $3}')
    echo -e "${GREEN}✅ Go 已安装: $GO_VERSION${NC}"
    
    # 检查 Go 版本
    GO_MAJOR=$(go version | grep -oP 'go\d+' | sed 's/go//' | cut -d. -f1)
    if [ "$GO_MAJOR" -lt 1 ]; then
        echo -e "${RED}❌ Go 版本过低，需要 Go 1.21+${NC}"
    fi
else
    echo -e "${RED}❌ Go 未安装${NC}"
fi
echo ""

# 检查 Node.js
echo -e "${BLUE}检查 Node.js 环境...${NC}"
if command -v node &> /dev/null; then
    NODE_VERSION=$(node --version)
    echo -e "${GREEN}✅ Node.js 已安装: $NODE_VERSION${NC}"
else
    echo -e "${RED}❌ Node.js 未安装${NC}"
fi
echo ""

# 检查 npm
echo -e "${BLUE}检查 npm...${NC}"
if command -v npm &> /dev/null; then
    NPM_VERSION=$(npm --version)
    echo -e "${GREEN}✅ npm 已安装: $NPM_VERSION${NC}"
else
    echo -e "${RED}❌ npm 未安装${NC}"
fi
echo ""

# 检查端口
echo -e "${BLUE}检查端口占用...${NC}"
if lsof -Pi :9000 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  端口 9000 已被占用${NC}"
    lsof -i :9000
else
    echo -e "${GREEN}✅ 端口 9000 可用${NC}"
fi

if lsof -Pi :9001 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo -e "${YELLOW}⚠️  端口 9001 已被占用${NC}"
    lsof -i :9001
else
    echo -e "${GREEN}✅ 端口 9001 可用${NC}"
fi
echo ""

# 检查后端依赖
echo -e "${BLUE}检查后端依赖...${NC}"
cd backend 2>/dev/null || {
    echo -e "${RED}❌ backend 目录不存在${NC}"
    exit 1
}

if [ -f "go.mod" ]; then
    echo -e "${GREEN}✅ go.mod 存在${NC}"
    if [ -f "go.sum" ]; then
        echo -e "${GREEN}✅ go.sum 存在${NC}"
    else
        echo -e "${YELLOW}⚠️  go.sum 不存在，需要运行 go mod download${NC}"
    fi
else
    echo -e "${RED}❌ go.mod 不存在${NC}"
fi
cd ..
echo ""

# 检查前端依赖
echo -e "${BLUE}检查前端依赖...${NC}"
cd frontend 2>/dev/null || {
    echo -e "${RED}❌ frontend 目录不存在${NC}"
    exit 1
}

if [ -f "package.json" ]; then
    echo -e "${GREEN}✅ package.json 存在${NC}"
    if [ -d "node_modules" ]; then
        echo -e "${GREEN}✅ node_modules 存在${NC}"
    else
        echo -e "${YELLOW}⚠️  node_modules 不存在，需要运行 npm install${NC}"
    fi
else
    echo -e "${RED}❌ package.json 不存在${NC}"
fi
cd ..
echo ""

# 检查日志文件
echo -e "${BLUE}检查日志文件...${NC}"
if [ -f "backend.log" ]; then
    echo -e "${YELLOW}📋 后端日志（最后10行）:${NC}"
    tail -10 backend.log
    echo ""
else
    echo -e "${GREEN}✅ 后端日志文件不存在（正常，服务未启动）${NC}"
fi

if [ -f "frontend.log" ]; then
    echo -e "${YELLOW}📋 前端日志（最后10行）:${NC}"
    tail -10 frontend.log
    echo ""
else
    echo -e "${GREEN}✅ 前端日志文件不存在（正常，服务未启动）${NC}"
fi

echo -e "${GREEN}✅ 诊断完成${NC}"
