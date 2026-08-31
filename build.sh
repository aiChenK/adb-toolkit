#!/usr/bin/env bash

set -e

# 如果存在 make 命令，直接委托给 Makefile
if command -v make >/dev/null 2>&1; then
    make build
    exit 0
fi

echo "🚀 [1/2] 构建前端静态资源..."
cd web
npm install
npm run build
cd ..

echo "🔨 [2/2] 编译 Go 单二进制应用程序..."
mkdir -p bin
go build -ldflags="-s -w" -o bin/adb-toolkit .

echo "✅ 构建成功！可执行文件路径: bin/adb-toolkit"
