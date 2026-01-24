# Memo Studio

一个简洁美观的笔记应用，支持 H5 和 Web 端，自适应设计，支持明暗主题切换。

## 技术栈

- **前端**: Svelte + Vite
- **后端**: Go + Gin + SQLite
- **特性**: 响应式设计、明暗主题、标签系统

## 项目结构

```
memo-studio/
├── backend/          # Go 后端服务
│   ├── main.go
│   ├── database/     # 数据库相关
│   ├── models/       # 数据模型
│   └── handlers/     # API 处理器
├── frontend/         # Svelte 前端应用
│   ├── src/
│   │   ├── components/  # 组件
│   │   ├── stores/      # 状态管理
│   │   └── utils/       # 工具函数
├── start.sh          # 一键启动脚本
└── README.md
```

## 快速开始

### 方式一：使用启动脚本（推荐）

```bash
./start.sh
```

### 方式二：手动启动

#### 后端启动

```bash
cd backend
go mod download
go run main.go
```

后端服务将在 `http://localhost:8080` 启动

#### 前端启动

```bash
cd frontend
npm install
npm run dev
```

前端应用将在 `http://localhost:3000` 启动

## 功能特性

### 第一版功能

- ✅ 笔记列表展示
- ✅ 笔记详情查看
- ✅ 新建笔记
- ✅ 标签支持
- ✅ 明暗主题切换
- ✅ 响应式设计（支持 H5 和 Web）

## API 接口

### 获取所有笔记
```
GET /api/notes
```

### 获取单个笔记
```
GET /api/notes/:id
```

### 创建笔记
```
POST /api/notes
Body: {
  "title": "笔记标题",
  "content": "笔记内容",
  "tags": ["标签1", "标签2"]
}
```

### 获取所有标签
```
GET /api/tags
```

## 开发说明

### 数据库

使用 SQLite 数据库，首次运行会自动创建数据库文件 `backend/notes.db` 和表结构。

数据库表结构：
- `users` - 用户表
- `notes` - 笔记表
- `tags` - 标签表
- `note_tags` - 笔记标签关联表

### API 接口

#### 笔记相关
- `GET /api/notes` - 获取所有笔记
- `GET /api/notes/:id` - 获取单个笔记
- `POST /api/notes` - 创建笔记
- `PUT /api/notes/:id` - 更新笔记
- `DELETE /api/notes/:id` - 删除笔记
- `DELETE /api/notes/batch` - 批量删除笔记

#### 标签相关
- `GET /api/tags` - 获取所有标签
- `PUT /api/tags/:id` - 更新标签
- `DELETE /api/tags/:id` - 删除标签
- `POST /api/tags/merge` - 合并标签

#### 认证相关（公开接口）
- `POST /api/auth/login` - 用户登录
- `POST /api/auth/register` - 用户注册

#### 用户相关（需要认证）
- `GET /api/auth/me` - 获取当前用户信息

### 主题切换

主题设置会保存在浏览器的 localStorage 中，支持明暗两种主题。

### 响应式设计

- 桌面端：网格布局，多列显示
- 移动端：单列布局，优化触摸交互
- iPad 端：2-3 列布局

## 用户认证

应用已实现用户认证功能：
- 用户注册：通过注册接口创建新账号
- 用户登录：使用用户名和密码登录
- JWT Token：登录后获取 token，用于后续 API 请求
- 路由保护：所有笔记和标签相关 API 都需要认证

### 首次使用

1. 启动后端服务
2. 打开前端应用
3. 点击"立即注册"创建账号
4. 注册成功后自动登录

## 后续计划

- [x] 笔记编辑功能
- [x] 笔记删除功能
- [x] 笔记搜索功能
- [x] 标签筛选功能
- [x] 用户认证
- [x] 数据导出
