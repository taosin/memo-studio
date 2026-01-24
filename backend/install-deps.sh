#!/bin/bash

# Go 依赖安装脚本

set -e

echo "📥 安装 Go 依赖..."

# 设置 Go 代理（使用国内镜像）
if [ -z "$GOPROXY" ]; then
    export GOPROXY=https://goproxy.cn,direct
    echo "📡 已设置 Go 代理为: goproxy.cn"
fi

# 进入后端目录
cd "$(dirname "$0")"

# 尝试安装依赖
echo "正在下载依赖..."
if go mod download; then
    echo "✅ 依赖下载成功"
else
    echo "⚠️  使用备用代理重试..."
    export GOPROXY=https://mirrors.aliyun.com/goproxy/,direct
    go mod download || {
        echo "❌ 依赖下载失败"
        echo ""
        echo "请尝试以下方法："
        echo "1. 检查网络连接"
        echo "2. 手动设置代理: export GOPROXY=https://goproxy.cn,direct"
        echo "3. 使用 VPN（如果需要）"
        exit 1
    }
fi

# 整理依赖
echo "正在整理依赖..."
go mod tidy

echo "✅ Go 依赖安装完成！"
echo ""
echo "现在可以运行: go run main.go"
