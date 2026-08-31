# 🤖 AGENTS.md - ADB Toolkit 开发者与 AI 协作规范

欢迎来到 **ADB Toolkit** 代码库。本文档旨在为 AI Agent 及人类开发者提供项目架构、设计决策、模块职责、接口契约与开发工作流的全局视图与操作准则。

---

## 1. 项目概述与设计哲学

- **项目定位**：轻量级、跨平台、本地运行的 Android 可视化调试与管理工具箱。
- **架构范式**：**前后端一体化单二进制**（Single Binary Distribution）。
  - 前端基于 Vue 3 + Ant Design Vue 开发；
  - 生产构建后通过 Go 1.16+ 的 `//go:embed web/dist/*` 机制内嵌至后端 Go 二进制；
  - 启动后作为一个独立的 Web 服务运行，自动唤起默认浏览器，零外部环境依赖（仅需系统安装有 `adb`）。

---

## 2. 技术栈清单

| 分层 / 领域 | 技术选型 | 版本 / 依赖说明 |
| :--- | :--- | :--- |
| **后端语言** | Go (Golang) | `>= 1.20` |
| **Web 框架** | Gin Web Framework | `github.com/gin-gonic/gin v1.9.1` |
| **静态内嵌** | Go Standard Library | `embed`, `io/fs`, `net/http` |
| **前端框架** | Vue 3 | `vue@^3.2.13` (Options API / Composition API) |
| **UI 组件库** | Ant Design Vue | `ant-design-vue@^4.0.0` |
| **HTTP 客户端** | Axios | `axios@^1.6.8` |
| **构建中心** | GNU Make / Shell | `Makefile`, `build.sh` |

---

## 3. 目录拓扑与模块职责

```text
adb-toolkit/
├── controllers/          # HTTP 请求控制器
│   ├── adbController.go  # 处理 ADB 指令执行与设备列表查询
│   └── toolController.go # 处理局域网 IP 查询等系统工具接口
├── core/                 # 核心业务逻辑
│   ├── adb.go            # ADB 进程调用、参数解析、设备扫描与重连
│   └── tool.go           # 跨平台网卡 IP 识别、ADB 检查、浏览器唤起
├── dto/                  # 数据传输对象
│   └── request/
│       └── Command.go    # ADB 命令表单请求定义 (CommandForm)
├── routes/               # 路由与中间件
│   ├── mid.go            # 静态内嵌文件系统 (EmbedFolder) 中间件
│   └── route.go          # API 路由注册与 CORS 跨域配置
├── web/                  # 前端 Vue 工程源码
│   ├── src/
│   │   ├── components/   # 组件 (AdbForm.vue 等)
│   │   ├── utils/        # Axios 封装 (http.js)
│   │   ├── App.vue       # 根组件
│   │   └── main.js       # 前端入口
│   ├── public/           # HTML 模板与公共图标
│   ├── vue.config.js     # 前端开发代理与构建配置
│   └── package.json
├── bin/                  # 编译产物目录 (已被 .gitignore 忽略)
├── main.go               # 应用启动入口、命令行参数解析、内嵌路由挂载
├── Makefile              # 核心构建与研发命令自动化
├── build.sh              # 便捷构建脚本 (兼容非 make 环境)
├── go.mod / go.sum       # Go 依赖清单
├── README.md             # 用户使用说明文档
└── AGENTS.md             # 本规范文档
```

---

## 4. 核心接口契约 (API Specification)

所有接口统一返回 JSON 格式：
```json
{
  "success": true,
  "errMessage": "",
  "data": ...
}
```

### 4.1 执行 ADB 命令
- **端点**：`POST /adb`
- **请求体 (form / json)**：
  ```json
  {
    "ip": "192.168.1.100",
    "port": "5555",
    "op": "setProxy | delProxy | clear | stop | free",
    "proxyAddr": "192.168.1.50:8888",
    "packageName": "com.example.app",
    "cmd": "shell getprop ro.product.model"
  }
  ```
- **参数说明**：
  - `op = free`：自由命令模式，直接执行 `cmd` 字段传入的命令（支持带引号参数）。
  - `op = setProxy`：设置指定设备的全局 HTTP 代理 (`proxyAddr`)。
  - `op = delProxy`：清除全局 HTTP 代理 (`:0`)。
  - `op = clear`：清理指定包名的应用缓存与数据。
  - `op = stop`：强制停止指定包名的应用进程。

### 4.2 获取已连接设备列表
- **端点**：`GET /devices`
- **响应示例**：
  ```json
  {
    "success": true,
    "errMessage": "",
    "data": [
      { "id": "192.168.1.100:5555", "status": "device", "type": "tcpip" },
      { "id": "emulator-5554", "status": "device", "type": "usb" }
    ]
  }
  ```

### 4.3 获取本机局域网 IP
- **端点**：`GET /ip`
- **响应示例**：
  ```json
  {
    "success": true,
    "errMessage": "",
    "data": "192.168.1.50"
  }
  ```

---

## 5. 常用开发与构建工作流

本项目所有研发指令均通过 `Makefile` 标准化：

```bash
# 1. 帮助查看
make help

# 2. 本地联调
make dev            # 一键同时启动前端与后端开发服务
# 或分终端运行：
make dev-backend    # 终端 A: 启动后端服务 (运行于 http://localhost:8088)
make dev-frontend   # 终端 B: 启动前端热重载 (运行于 http://localhost:8080，代理转发至 :8088)

# 3. 生产构建 (前端编译 + 单二进制内嵌打包)
make                # 编译当前宿主机平台的二进制到 bin/adb-toolkit

# 4. 交叉编译多平台发布包 (支持 VERSION=v1.0.0 或 V=v1.0.0)
make build-darwin VERSION=v1.0.0   # 编译 macOS arm64 + amd64
make build-linux VERSION=v1.0.0    # 编译 Linux amd64
make build-windows VERSION=v1.0.0  # 编译 Windows amd64 (.exe)
make build-all VERSION=v1.0.0      # 一键编译全平台包 (输出形如 adb-toolkit-v1.0.0-darwin-amd64)

# 5. 清理构建缓存
make clean          # 删除 bin/ 与 web/dist/
```

---

## 6. Agent 编码与修改准则

1. **内嵌产物生命周期**：
   - 禁止手动修改 `web/dist` 中的静态文件，所有前端修改必须在 `web/src/` 中进行，并通过 `npm run build` 生成最新产物。
   - 在测试 Go 二进制打包前，必须先确保 `web/dist` 处于最新状态。
2. **命令执行健壮性与错误处理**：
   - 执行 ADB 命令必须使用 `parseCommandArgs` 或切片化参数，严禁直接使用简单空格切分导致破坏复杂参数。
   - 命令执行必须捕获 stdout 与 stderr，即使执行退出码非 0，也必须将具体错误内容返回给前端控制台展示。
3. **定向设备操作约束**：
   - 设备状态异常重连或断开时，严禁使用全局 `adb disconnect`，必须使用 `adb disconnect <target>` 定向断开，防止误断开开发者其他在线设备。
4. **Git 与依赖管理**：
   - 严禁将大体积二进制文件、平台编译包或 node_modules 提交至 Git 仓库。
   - 所有新增环境忽略项应及时补充至 `.gitignore`。
