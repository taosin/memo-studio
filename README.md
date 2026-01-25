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

## 新一代实现（Go + SQLite FTS5 + SvelteKit）

你这次要求的「Go 后端托管/内嵌 SvelteKit 静态文件」对应如下脚本：

- **开发模式（热更新）**：

```bash
./dev-kit.sh
```

打开 `http://localhost:9001`（SvelteKit dev），API 走代理到 `http://localhost:9000`。

- **生产构建 + 一键启动（Go 直接提供前端静态文件）**：

```bash
./start-prod.sh
```

启动后直接打开 `http://localhost:9000`。

说明：
- Go 构建时启用了 `sqlite_fts5` build tag（用于 SQLite FTS5）
- SvelteKit 构建产物会同步到 `backend/public/`，并由 Go 在运行时托管（SPA fallback 到 `index.html`）

## Docker 发布/自部署（推荐）

目标：别人可以 `docker run` 或 `docker compose up` 在 NAS/服务器上自部署，并且数据可持久化。

### 1) 最小可运行（docker run）

必填环境变量：
- **`MEMO_JWT_SECRET`**：JWT 密钥（生产必须设置，建议 32+ 字符）

推荐环境变量：
- **`MEMO_ADMIN_PASSWORD`**：用于初始化/重置管理员 `admin` 的密码（不设置则首次启动会随机生成并打印到容器日志）
- **`MEMO_CORS_ORIGINS`**：允许的前端域名（逗号分隔），例如 `https://your.domain,https://nas.local`

容器数据目录：
- **`/data/notes.db`**：SQLite 数据库（通过 `MEMO_DB_PATH` 指定）
- **`/data/storage`**：附件目录（通过 `MEMO_STORAGE_DIR` 指定）

示例：

```bash
docker run -d \
  --name memo-studio \
  -p 9000:9000 \
  -v memo_data:/data \
  -e MEMO_JWT_SECRET="please-change-me" \
  -e MEMO_ADMIN_PASSWORD="your-strong-password" \
  -e MEMO_ENV=production \
  -e GIN_MODE=release \
  memo-studio:local
```

启动后访问：`http://localhost:9000`

默认管理员用户名：`admin`

### 2) docker compose（最推荐）

直接使用仓库根目录的 `docker-compose.yml`：

```bash
docker compose up -d --build
```

### 3) 首次管理员策略（重要）

- 生产环境**不再固定** `admin/admin123`
- 若设置了 `MEMO_ADMIN_PASSWORD`：启动时会确保 `admin` 存在并重置密码，同时标记“需要修改密码”
- 若没设置且数据库为空：启动会生成随机初始密码，并打印到日志（请登录后立即修改）

### 4) 生产配置（环境变量）

- **`PORT`**：监听端口（默认 9000）
- **`MEMO_DB_PATH`**：SQLite 路径（默认 `./notes.db`；容器建议 `/data/notes.db`）
- **`MEMO_STORAGE_DIR`**：附件目录（默认 `./storage`；容器建议 `/data/storage`）
- **`MEMO_CORS_ORIGINS`**：CORS 白名单（逗号分隔；不填默认放开）
- **`MEMO_JWT_SECRET`**：JWT 密钥（生产必须设置）

### 5) 多架构镜像（NAS 兼容）

建议发布 `linux/amd64` 与 `linux/arm64` 两种架构镜像（群晖/威联通/树莓派常用）。
后续可以用 GitHub Actions + buildx 自动构建并推送到 Docker Hub/GHCR。

## 镜像发布（给别人 docker pull）

本仓库已内置 GitHub Actions：推送 tag（如 `v0.1.0`）会自动构建并推送镜像到 GHCR：

- 镜像地址：`ghcr.io/<你的GitHub用户名>/<仓库名>:latest`
- 也会推送版本 tag：例如 `ghcr.io/<你的GitHub用户名>/<仓库名>:v0.1.0`

### 发布步骤

```bash
git tag v0.1.0
git push origin v0.1.0
```

Actions 运行完成后，别人即可部署：

```bash
docker run -d \
  --name memo-studio \
  -p 9000:9000 \
  -v memo_data:/data \
  -e MEMO_JWT_SECRET="please-change-me" \
  -e MEMO_ADMIN_PASSWORD="your-strong-password" \
  ghcr.io/<你的GitHub用户名>/<仓库名>:latest
```

## 同步发布到 Docker Hub（可选）

很多 NAS 更习惯从 Docker Hub 拉取镜像。本仓库工作流已支持在打 tag 时**同步推送 Docker Hub**，前提是你配置好 Secrets：

在 GitHub 仓库 `Settings → Secrets and variables → Actions` 添加：
- **`DOCKERHUB_USERNAME`**：你的 Docker Hub 用户名
- **`DOCKERHUB_TOKEN`**：Docker Hub Access Token（建议用 token，不要用密码）

发布后 Docker Hub 镜像名规则：
- `docker.io/<DOCKERHUB_USERNAME>/<仓库名>:latest`
- `docker.io/<DOCKERHUB_USERNAME>/<仓库名>:v0.1.0`

别人部署示例：

```bash
docker run -d \
  --name memo-studio \
  -p 9000:9000 \
  -v memo_data:/data \
  -e MEMO_JWT_SECRET="please-change-me" \
  -e MEMO_ADMIN_PASSWORD="your-strong-password" \
  docker.io/<你的DockerHub用户名>/<仓库名>:latest
```

## 自建 AI CR 机器人（PR 自动审查）

本仓库已内置工作流：`.github/workflows/ai-pr-review.yml`  
默认不会对所有 PR 自动评论（避免刷屏/与 Gemini 等机器人冲突）。运行方式如下：

- **方式 A（推荐）**：给 PR 加上标签 **`ai-review`**，工作流就会自动运行并更新同一条评论
- **方式 B**：手动触发 `AI PR Review` 工作流，并填写 `pr_number`

工作流会：
- 拉取 PR diff（不 checkout PR 分支代码，避免安全风险）
- 调用你配置的大模型 API
- 在 PR 下发布/更新一条中文 CR 评论

### 配置 Secrets

在 GitHub 仓库：`Settings → Secrets and variables → Actions → New repository secret` 添加：

- **`AI_REVIEW_API_KEY`**：模型 API Key（必填）
- **`AI_REVIEW_MODEL`**：模型名（必填，例如 `gpt-4o-mini` / `gpt-4.1-mini` / 你自建模型名）
- **`AI_REVIEW_BASE_URL`**：可选，OpenAI 兼容接口地址（默认 `https://api.openai.com/v1`）

说明：
- 如果未配置 `AI_REVIEW_API_KEY/AI_REVIEW_MODEL`，工作流会自动跳过（不报错）

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

### 健康检查（公开接口）
- `GET /health` - 健康检查端点

#### 认证相关（公开接口）
- `POST /api/auth/login` - 用户登录
  - 请求体: `{ "username": "string", "password": "string" }`
  - 返回: `{ "token": "string", "user": {...} }`
- `POST /api/auth/register` - 用户注册
  - 请求体: `{ "username": "string", "password": "string", "email": "string" }`
  - 返回: `{ "token": "string", "user": {...} }`

#### 用户相关（需要认证）
- `GET /api/auth/me` - 获取当前用户信息
  - 需要 Authorization: Bearer <token>
  - 返回: `{ "id": number, "username": "string", "email": "string", "created_at": "datetime" }`

### 笔记相关（需要认证）

- `GET /api/notes` - 获取所有笔记
  - 返回: `[{ "id": number, "title": "string", "content": "string", "tags": [...], "created_at": "datetime", "updated_at": "datetime" }]`
- `GET /api/notes/:id` - 获取单个笔记
  - 返回: `{ "id": number, "title": "string", "content": "string", "tags": [...], "created_at": "datetime", "updated_at": "datetime" }`
- `POST /api/notes` - 创建笔记
  - 请求体: `{ "title": "string", "content": "string", "tags": ["string"] }`
  - 返回: 创建的笔记对象
- `PUT /api/notes/:id` - 更新笔记
  - 请求体: `{ "title": "string", "content": "string", "tags": ["string"] }`
  - 返回: 更新后的笔记对象
- `DELETE /api/notes/:id` - 删除笔记
  - 返回: `{ "success": true, "message": "笔记已删除" }`
- `DELETE /api/notes/batch` - 批量删除笔记
  - 请求体: `{ "ids": [number] }`
  - 返回: `{ "success": true, "deleted": number, "message": "string" }`

### 标签相关（需要认证）

- `GET /api/tags` - 获取所有标签
  - 返回: `[{ "id": number, "name": "string", "color": "string", "created_at": "datetime" }]`
- `PUT /api/tags/:id` - 更新标签
  - 请求体: `{ "name": "string", "color": "string" }`
  - 返回: 更新后的标签对象
- `DELETE /api/tags/:id` - 删除标签
  - 返回: `{ "success": true, "message": "标签已删除" }`
- `POST /api/tags/merge` - 合并标签
  - 请求体: `{ "sourceId": number, "targetId": number }`
  - 返回: `{ "success": true, "message": "标签合并成功" }`

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

### 热更新说明

#### 前端热更新（自动）✅

前端使用 Vite，**默认支持热模块替换（HMR）**：
- ✅ 修改前端代码后，浏览器会自动刷新
- ✅ 无需手动重启前端服务
- ✅ 修改样式和组件会立即生效
- ✅ 保持应用状态（不会丢失数据）

**使用方式：**
1. 启动服务后，修改 `frontend/src/` 下的任何文件
2. 保存文件后，浏览器会自动更新
3. 无需任何手动操作

#### 后端热重载（需要工具）

后端 Go 服务默认不支持热重载，有两种方式：

**方式一：使用 Air（推荐，自动热重载）**

1. 安装 Air：
   ```bash
   go install github.com/cosmtrek/air@latest
   ```

2. 在 `backend` 目录运行：
   ```bash
   cd backend
   ./start-dev.sh
   # 或直接运行
   air
   ```

3. 修改 Go 代码后，Air 会自动重新编译和重启服务

**方式二：手动重启（简单但需要手动操作）**

修改代码后，需要手动停止并重新启动后端服务：
```bash
# 停止服务（Ctrl+C）
# 然后重新运行
cd backend
go run main.go
```

### 日志文件

- `backend.log` - 后端服务日志
- `frontend.log` - 前端开发服务器日志

### 停止服务

**macOS / Linux:**
在运行 `./start.sh` 的终端中按 `Ctrl+C` 即可停止所有服务。

**Windows:**
关闭运行 `start.bat` 的命令窗口即可停止所有服务。

## 故障排查

### 端口被占用

如果 9000 或 9001 端口被占用，启动脚本会自动尝试清理。如果失败，请手动停止占用端口的进程：

```bash
# 查看端口占用
lsof -i :9000
lsof -i :9001

# 停止进程（替换 PID）
kill -9 <PID>
```

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

### 热更新不工作

**前端：**
- 检查浏览器控制台是否有错误
- 尝试硬刷新（Ctrl+Shift+R 或 Cmd+Shift+R）
- 检查 Vite 开发服务器是否正常运行

**后端：**
- 确保已安装 Air：`go install github.com/cosmtrek/air@latest`
- 检查 `.air.toml` 配置文件是否存在
- 查看 Air 的输出日志

## 许可证

MIT License
