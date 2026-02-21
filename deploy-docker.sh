#!/bin/bash
# ============================================
# Memo Studio 一键 Docker 部署脚本
# ============================================

set -e

# 颜色
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
RED='\033[0;31m'
NC='\033[0m'

SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
cd "$SCRIPT_DIR"

echo -e "${BLUE}🚀 Memo Studio Docker 一键部署${NC}"
echo "================================"

# 检查 Docker
if ! command -v docker &> /dev/null; then
    echo -e "${RED}❌ Docker 未安装，请先安装 Docker${NC}"
    exit 1
fi

if ! command -v docker compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo -e "${RED}❌ Docker Compose 未安装${NC}"
    exit 1
fi

# 检查 .env 文件
if [ ! -f ".env" ]; then
    echo -e "${YELLOW}📝 未发现 .env 文件，创建中...${NC}"
    if [ -f ".env.example" ]; then
        cp .env.example .env
        echo -e "${YELLOW}⚠️  请编辑 .env 文件设置 MEMO_JWT_SECRET！${NC}"
        exit 1
    else
        # 自动生成 JWT Secret
        JWT_SECRET=$(openssl rand -base64 32 2>/dev/null || head -c 32 /dev/urandom | base64)
        cat > .env << EOF
MEMO_JWT_SECRET=$JWT_SECRET
MEMO_ADMIN_PASSWORD=
MEMO_CORS_ORIGINS=
EOF
        echo -e "${GREEN}✅ 已自动生成 JWT Secret${NC}"
    fi
fi

# 提示 JWT Secret
if grep -q "please-change" .env 2>/dev/null; then
    echo -e "${RED}⚠️  请编辑 .env 文件，将 MEMO_JWT_SECRET 设置为强密码！${NC}"
    exit 1
fi

echo -e "${GREEN}✅ 环境配置检查通过${NC}"

# 停止旧容器
echo -e "${YELLOW}🛑 停止旧容器...${NC}"
docker compose down 2>/dev/null || true

# 构建并启动
echo -e "${BLUE}🔨 构建镜像...${NC}"
docker compose build --no-cache

echo -e "${BLUE}🚀 启动服务...${NC}"
docker compose up -d

# 等待健康检查
echo -e "${YELLOW}⏳ 等待服务启动...${NC}"
for i in {1..30}; do
    if curl -s -f "http://localhost:9000/health" >/dev/null 2>&1; then
        echo -e "${GREEN}✅ 服务已就绪！${NC}"
        break
    fi
    if [ $i -eq 30 ]; then
        echo -e "${RED}❌ 服务启动超时，请检查日志：docker compose logs${NC}"
        exit 1
    fi
    sleep 1
done

echo ""
echo -e "${GREEN}🎉 部署成功！${NC}"
echo "================================"
echo -e "${BLUE}🌐 访问地址: ${GREEN}http://localhost:9000${NC}"
echo -e "${YELLOW}📝 查看日志: ${GREEN}docker compose logs -f${NC}"
echo -e "${YELLOW}🛑 停止服务: ${GREEN}docker compose down${NC}"
echo ""
echo -e "${YELLOW}💡 提示：生产环境请配置域名和 HTTPS${NC}"
