#!/bin/bash

# 清理端口占用脚本

echo "🔍 检查端口占用情况..."

# 检查 9000 端口
if lsof -Pi :9000 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "端口 9000 被占用:"
    lsof -i :9000
    read -p "是否停止占用 9000 端口的进程? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        lsof -ti:9000 | xargs kill -9 2>/dev/null
        echo "✅ 已停止占用 9000 端口的进程"
    fi
else
    echo "✅ 端口 9000 未被占用"
fi

# 检查 9001 端口
if lsof -Pi :9001 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "端口 9001 被占用:"
    lsof -i :9001
    read -p "是否停止占用 9001 端口的进程? (y/n) " -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        lsof -ti:9001 | xargs kill -9 2>/dev/null
        echo "✅ 已停止占用 9001 端口的进程"
    fi
else
    echo "✅ 端口 9001 未被占用"
fi
