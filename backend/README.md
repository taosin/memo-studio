# Memo Studio Backend

基于 Go + Gin 的笔记应用后端服务

## 依赖

- Go 1.21+
- Gin Web Framework
- SQLite3

## 安装依赖

```bash
go mod download
```

## 运行

```bash
go run main.go
```

服务将在 `http://localhost:9000` 启动

## API 文档

### 笔记相关

- `GET /api/notes` - 获取所有笔记
- `GET /api/notes/:id` - 获取单个笔记
- `POST /api/notes` - 创建笔记

### 标签相关

- `GET /api/tags` - 获取所有标签

## 数据库

使用 SQLite 数据库，数据库文件为 `notes.db`，首次运行会自动创建表结构。
