# Makefile for ADB Toolkit

APP_NAME := adb-toolkit
BIN_DIR := bin
WEB_DIR := web
LDFLAGS := -s -w

.PHONY: help all build build-web build-app build-darwin build-darwin-arm64 build-darwin-amd64 build-linux build-windows build-all run-backend run-frontend dev dev-backend dev-frontend clean test

# Default target
all: build

help: ## 显示帮助信息
	@echo "🛠️  ADB Toolkit 构建工具"
	@echo ""
	@echo "用法: make [目标]"
	@echo ""
	@echo "常用目标:"
	@echo "  make dev               - 一键同时启动前端与后端开发服务"
	@echo "  make dev-backend       - 本地运行后端服务 (端口 8088)"
	@echo "  make dev-frontend      - 本地启动前端开发热重载服务 (端口 8080)"
	@echo "  make build             - 构建前端静态资源并编译当前平台的单二进制文件"
	@echo "  make build-web         - 仅构建前端静态产物 (web/dist)"
	@echo "  make build-app         - 仅编译当前平台的 Go 二进制文件"
	@echo "  make build-darwin      - 编译 macOS 通用/当前架构二进制文件"
	@echo "  make build-linux       - 交叉编译 Linux (amd64) 二进制文件"
	@echo "  make build-windows     - 交叉编译 Windows (amd64) 可执行文件"
	@echo "  make build-all         - 一键构建多平台发行包 (macOS, Linux, Windows)"
	@echo "  make clean             - 清理构建产物 (bin/ 及 web/dist/)"
	@echo "  make test              - 运行后端单元测试"
	@echo ""

build: build-web build-app ## 构建前端并编译当前平台二进制

build-web: ## 构建前端
	@echo "📦 [1/2] 构建前端静态资源..."
	@cd $(WEB_DIR) && npm install && npm run build

build-app: ## 编译当前系统架构的 Go 二进制
	@echo "🔨 [2/2] 编译 Go 应用程序..."
	@mkdir -p $(BIN_DIR)
	@go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(APP_NAME) .
	@echo "✅ 构建完成: $(BIN_DIR)/$(APP_NAME)"

build-darwin-arm64: build-web ## 编译 macOS Apple Silicon (arm64)
	@echo "🍎 交叉编译 macOS (arm64)..."
	@mkdir -p $(BIN_DIR)
	@GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(APP_NAME)-darwin-arm64 .

build-darwin-amd64: build-web ## 编译 macOS Intel (amd64)
	@echo "🍎 交叉编译 macOS (amd64)..."
	@mkdir -p $(BIN_DIR)
	@GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(APP_NAME)-darwin-amd64 .

build-darwin: build-darwin-arm64 build-darwin-amd64 ## 编译 macOS 双架构

build-linux: build-web ## 交叉编译 Linux (amd64)
	@echo "🐧 交叉编译 Linux (amd64)..."
	@mkdir -p $(BIN_DIR)
	@GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(APP_NAME)-linux-amd64 .

build-windows: build-web ## 交叉编译 Windows (amd64)
	@echo "🪟 交叉编译 Windows (amd64)..."
	@mkdir -p $(BIN_DIR)
	@GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(APP_NAME)-windows-amd64.exe .

build-all: build-web ## 一键交叉编译多平台版本
	@echo "🚀 开始全平台交叉编译..."
	@mkdir -p $(BIN_DIR)
	@GOOS=darwin GOARCH=arm64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(APP_NAME)-darwin-arm64 .
	@GOOS=darwin GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(APP_NAME)-darwin-amd64 .
	@GOOS=linux GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(APP_NAME)-linux-amd64 .
	@GOOS=windows GOARCH=amd64 go build -ldflags="$(LDFLAGS)" -o $(BIN_DIR)/$(APP_NAME)-windows-amd64.exe .
	@echo "🎉 全平台构建成功！产物位于 $(BIN_DIR)/"

dev: ## 同时启动前后端开发服务
	@echo "🚀 正在同时启动前端与后端开发服务..."
	@trap 'kill 0' SIGINT SIGTERM EXIT; \
	go run . & \
	(cd $(WEB_DIR) && npm run serve) & \
	wait

dev-backend: ## 运行后端 (开发环境)
	@echo "🚀 启动后端开发服务..."
	@go run .

run-backend: dev-backend

dev-frontend: ## 运行前端 (开发环境)
	@echo "🚀 启动前端开发服务..."
	@cd $(WEB_DIR) && npm run serve

run-frontend: dev-frontend

test: ## 运行测试
	@echo "🧪 运行 Go 测试..."
	@go test -v ./...

clean: ## 清理产物
	@echo "🧹 清理构建产物..."
	@rm -rf $(BIN_DIR)/ $(WEB_DIR)/dist/
	@echo "✨ 清理完成"
