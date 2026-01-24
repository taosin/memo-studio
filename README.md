# Memo Studio

一个简洁美观的笔记应用，支持 H5 和 Web 端，自适应设计，支持明暗主题切换。

## 技术栈

- **前端**: Svelte + Vite
- **后端**: Go + Gin + SQLite
- **特性**: 响应式设计、明暗主题、标签系统、用户认证

## 快速开始

### 一键启动（推荐）

**macOS / Linux:**
```bash
./start.sh
```

**Windows:**
```cmd
start.bat
```

脚本会自动：
- ✅ 检查 Go 和 Node.js 环境
- ✅ 安装依赖（Go modules 和 npm packages）
- ✅ 检查并清理端口占用
- ✅ 启动后端服务（:9000）
- ✅ 启动前端应用（:9001）
- ✅ 等待服务就绪后显示访问地址

启动成功后：
- 📝 后端 API: http://localhost:9000
- 🌐 前端应用: http://localhost:9001

### 手动启动

#### 1. 启动后端

```bash
cd backend
go mod download
go run main.go
```

后端服务将在 `http://localhost:9000` 启动

#### 2. 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端应用将在 `http://localhost:9001` 启动

## 首次使用

1. 启动服务后，打开浏览器访问 http://localhost:9001
2. 点击"立即注册"创建账号
3. 注册成功后自动登录，即可开始使用

## 项目结构

```
memo-studio/
├── backend/          # Go 后端服务
│   ├── main.go       # 入口文件
│   ├── database/     # 数据库相关
│   ├── models/       # 数据模型
│   ├── handlers/     # API 处理器
│   ├── middleware/   # 中间件
│   └── utils/        # 工具函数
├── frontend/         # Svelte 前端应用
│   ├── src/
│   │   ├── components/  # 组件
│   │   ├── stores/       # 状态管理
│   │   └── utils/       # 工具函数
│   └── vite.config.js
├── start.sh          # 一键启动脚本
└── README.md
```

## 功能特性

### 已实现功能

- ✅ 用户注册和登录（JWT 认证）
- ✅ 笔记列表展示（瀑布流/时间线模式）
- ✅ 笔记详情查看
- ✅ 新建/编辑笔记
- ✅ 删除笔记（单个/批量）
- ✅ 标签系统（创建、编辑、删除、合并）
- ✅ 高级搜索（关键词、日期、标签筛选）
- ✅ 数据导出（Markdown、JSON、CSV）
- ✅ 明暗主题切换
- ✅ 响应式设计（支持 H5 和 Web）
- ✅ 热力图显示

## API 接口

### 认证相关（公开接口）

- `POST /api/auth/login` - 用户登录
- `POST /api/auth/register` - 用户注册

### 用户相关（需要认证）

- `GET /api/auth/me` - 获取当前用户信息

### 笔记相关（需要认证）

- `GET /api/notes` - 获取所有笔记
- `GET /api/notes/:id` - 获取单个笔记
- `POST /api/notes` - 创建笔记
- `PUT /api/notes/:id` - 更新笔记
- `DELETE /api/notes/:id` - 删除笔记
- `DELETE /api/notes/batch` - 批量删除笔记

### 标签相关（需要认证）

- `GET /api/tags` - 获取所有标签
- `PUT /api/tags/:id` - 更新标签
- `DELETE /api/tags/:id` - 删除标签
- `POST /api/tags/merge` - 合并标签

## 数据库

使用 SQLite 数据库，首次运行会自动创建数据库文件 `backend/notes.db` 和表结构。

数据库表结构：
- `users` - 用户表
- `notes` - 笔记表
- `tags` - 标签表
- `note_tags` - 笔记标签关联表

## 开发说明

### 环境要求

- Go 1.21+
- Node.js 18+
- npm 或 yarn

### 开发模式

后端和前端都支持热重载，修改代码后会自动重新编译。

### 日志文件

- `backend.log` - 后端服务日志
- `frontend.log` - 前端开发服务器日志

### 停止服务

**macOS / Linux:**
在运行 `./start.sh` 的终端中按 `Ctrl+C` 即可停止所有服务。

**Windows:**
关闭运行 `start.bat` 的命令窗口即可停止所有服务。

## 故障排查

### 使用诊断工具

如果遇到启动问题，可以先运行诊断脚本：

```bash
./check.sh
```

诊断脚本会检查：
- Go 和 Node.js 环境
- 端口占用情况
- 依赖安装状态
- 日志文件内容

### 查看错误日志

启动脚本会在项目根目录生成日志文件：

```bash
# 查看后端日志
tail -f backend.log

# 查看前端日志
tail -f frontend.log
```

### 常见问题

#### 1. 端口被占用

如果 9000 或 9001 端口被占用，启动脚本会自动尝试清理。如果失败，请手动停止占用端口的进程：

```bash
# 查看端口占用
lsof -i :9000
lsof -i :9001

# 停止进程（替换 PID）
kill -9 <PID>
```

#### 2. Go 依赖安装失败

如果 Go 依赖安装失败，可能是网络问题：

```bash
cd backend
go mod download
go mod tidy
```

#### 3. npm 依赖安装失败

如果 npm 依赖安装失败：

```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

#### 4. 后端启动失败

检查后端日志文件 `backend.log`，常见原因：
- 数据库文件权限问题
- 端口被占用
- Go 依赖缺失

#### 5. 前端启动失败

检查前端日志文件 `frontend.log`，常见原因：
- npm 依赖缺失
- 端口被占用
- Vite 配置问题

### 依赖安装失败

**Go 依赖：**
```bash
cd backend
go mod download
go mod tidy
```

**npm 依赖：**
```bash
cd frontend
rm -rf node_modules package-lock.json
npm install
```

### 数据库问题

如果数据库文件损坏，可以删除后重新启动：

```bash
cd backend
rm notes.db
# 重新启动服务，数据库会自动创建
```

## 许可证

MIT License
