#!/bin/bash

# 后端开发模式启动脚本（使用 Air 热重载）

set -e

# 颜色定义
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

echo -e "${GREEN}🚀 启动后端开发服务器（热重载模式）...${NC}"

# 检查 Air 是否安装
if ! command -v air &> /dev/null; then
    echo -e "${YELLOW}⚠️  Air 未安装，正在安装...${NC}"
    go install github.com/cosmtrek/air@latest
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Air 安装失败${NC}"
        echo -e "${YELLOW}   请手动安装: go install github.com/cosmtrek/air@latest${NC}"
        exit 1
    fi
    echo -e "${GREEN}✅ Air 安装成功${NC}"
fi

# 设置 Go 代理（如果需要）
if [ -z "$GOPROXY" ]; then
    export GOPROXY=https://goproxy.cn,direct
    echo -e "${YELLOW}📡 已设置 Go 代理为: goproxy.cn${NC}"
fi

# 检查依赖
if [ ! -f "go.sum" ]; then
    echo -e "${YELLOW}📥 安装 Go 依赖...${NC}"
    go mod download && go mod tidy
fi

# 启动 Air（热重载）
echo -e "${GREEN}✅ 后端服务已启动（热重载模式）${NC}"
echo -e "${YELLOW}💡 修改代码后会自动重新编译和重启${NC}"
echo ""

air
