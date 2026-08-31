# 🛠️ ADB Toolkit (ADB 可视化调试工具箱)

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.20+-00ADD8?style=flat&logo=go" alt="golang" />
  <img src="https://img.shields.io/badge/Vue-3.x-4FC08D?style=flat&logo=vue.js" alt="vue" />
  <img src="https://img.shields.io/badge/Ant%20Design%20Vue-4.x-1890FF?style=flat&logo=antdesign" alt="antd" />
  <img src="https://img.shields.io/badge/License-MIT-green" alt="license" />
</p>

**ADB Toolkit** 是一个基于 **Go (Gin) + Vue 3 (Ant Design Vue)** 构建的轻量级、跨平台 Android 可视化调试工具箱。

通过现代化 Web 界面简化常用的 ADB 操作，支持多设备/客户端分组独立管理、已连接设备自动扫描、局域网抓包代理一键配置、应用进程/缓存清理与自由命令执行。内置静态资源打包机制，编译后仅为一个**绿色无依赖的单可执行文件**，启动后自动唤起浏览器使用。

---

## ✨ 核心特性

- 🖥️ **Web 可视化交互**：告别繁琐记忆 ADB 命令行，界面简洁直观，启动自动拉起浏览器。
- 📱 **多设备与分组管理**：支持添加多个设备配置卡片，支持自动扫描当前连接的 USB/无线 Android 设备。
- 🌐 **一键代理配置**：
  - 智能跨平台获取当前开发机（macOS / Linux / Windows）的局域网物理 IPv4 地址。
  - 一键为指定 Android 设备设置/清除全局 HTTP 代理（适合抓包与联调）。
- ⚡ **应用快捷操作**：一键清除 App 缓存（`pm clear`）、强制停止 App 进程（`am force-stop`）。
- 📝 **命令控制台与预置命令**：
  - 支持常用预置命令快捷选择（`devices`、`disconnect`、`pm list packages -3`、`dumpsys battery` 等）。
  - 支持带引号/参数的自由命令输入并指定设备定向执行。
  - 实时色彩化时间线输出执行结果与日志，支持单条日志一键复制。
- 📦 **单二进制便携分发**：前端静态资源通过 Go `embed` 完整内嵌，无需额外配置 Nginx 或 Node 环境。

---

## 🏗️ 目录结构

```text
adb-toolkit/
├── web/                  # 前端工程 (Vue 3 + Ant Design Vue)
│   ├── src/              # 前端源码
│   ├── package.json
│   └── vue.config.js
├── controllers/          # Go 接口控制器
├── core/                 # ADB 执行核心与网络工具
├── dto/                  # 请求/响应数据结构
├── routes/               # Gin 路由与中间件定义
├── main.go               # 主入口（内嵌 web/dist）
├── go.mod                # Go 模块定义
├── Makefile              # 跨平台标准构建命令集
├── build.sh              # 一键打包脚本
├── AGENTS.md             # 架构设计与 AI Agent 开发指引
└── README.md             # 项目说明文档
```

---

## 🚀 快速上手

### 1. 前置依赖
- [Go](https://go.dev/) (>= 1.20)
- [Node.js](https://nodejs.org/) (>= 16) & npm
- 本机已配置并可正常运行 `adb` 环境变量

---

### 2. 本地开发与联调

#### 启动后端服务：
```bash
make dev-backend
# 或直接运行: go run .
```
后端服务默认监听 `http://localhost:8088`。支持参数自定义：`go run . -port=9000 -no-browser=true`。

#### 启动前端（热重载开发）：
```bash
make dev-frontend
# 或进入 web 目录: cd web && npm run serve
```
前端默认运行在 `http://localhost:8080`，已配置代理转发至后端的 `:8088` 端口。

---

### 3. 一键编译打包（生产环境）

通过 `Makefile` 构建前端并将静态产物内嵌到 Go 二进制中：

```bash
# 1. 编译当前平台单二进制文件 (输出至 bin/adb-toolkit)
make

# 2. 交叉编译 macOS (Apple Silicon + Intel)
make build-darwin

# 3. 交叉编译 Linux (amd64)
make build-linux

# 4. 交叉编译 Windows (amd64)
make build-windows

# 5. 一键编译全平台发行包
make build-all
```

编译完成后，只需将 `bin/` 下生成的单个可执行文件分发给使用者即可，双击直接运行！

---

## 📝 更新日志

### v1.1.0
- 🧹 彻底移除废弃的 Docker 镜像配置与大文件残留，仓库显著瘦身。
- 🛠️ 重构构建体系为现代 `Makefile`，支持全平台交叉编译（macOS arm64/amd64、Linux、Windows）。
- ⚡ 增强 ADB 命令参数解析与 stderr 错误捕获，支持定向设备断开。
- 🌐 优化跨平台局域网 IP 扫描算法，精准识别物理无线/有线网卡。
- 🎨 升级 Web 控制台体验：新增在线设备检测与快速填充、多预置命令、单条日志一键复制。
- 🤖 新增 `AGENTS.md` 规范化协作文档。
