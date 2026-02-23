# Docker 部署检查清单

## ✅ 已确认正常

- [x] 多阶段构建（分离 kit 和 go 构建）
- [x] 使用非 root 用户
- [x] 健康检查配置
- [x] 数据卷持久化
- [x] SQLite FTS5 支持
- [x] CORS 配置支持
- [x] API 两套路由兼容（/api/* 和 /api/v1/*）

## ⚠️ 已修复

### 1. JWT Secret 必须修改

**当前：**
```yaml
MEMO_JWT_SECRET: "change-me"
```

**问题：** 使用默认字符串，极其不安全

**修复：**
```bash
# 部署时必须设置强密码
export MEMO_JWT_SECRET=$(openssl rand -base64 32)
```

或者使用 docker-compose：
```yaml
MEMO_JWT_SECRET: "${JWT_SECRET:-$(openssl rand -base64 32)}"
```

### 2. API 路由版本不一致 ⚠️（已解决）

**代码定义：** `/api/v1/auth/login`
**实际调用：** `/api/auth/login`

**状态：** ✅ 已兼容
- 后端同时提供 `/api/v1/*`（新）和 `/api/*`（旧/兼容）两套 API
- 前端使用 `/api/*` 路由，能正常工作
- 日志显示请求正常匹配

### 3. docker-compose 版本属性过时

```yaml
version: "3.8"  # ❌ 已过时
```

**修复：** 直接删除 `version` 属性。

### 3. GIN_MODE 配置 ✅

**代码逻辑：**
```go
if strings.TrimSpace(os.Getenv("GIN_MODE")) == "" && strings.TrimSpace(os.Getenv("MEMO_ENV")) == "production" {
    gin.SetMode(gin.ReleaseMode)
}
```

**状态：** ✅ 已正确配置
- `GIN_MODE=release` 已设置
- `MEMO_ENV=production` 已设置
- 生产模式会自动启用

## 📋 部署前必做清单

### 部署命令

```bash
# 1. 生成强 JWT Secret
export JWT_SECRET=$(openssl rand -base64 32)
echo "JWT_SECRET=$JWT_SECRET"

# 2. 设置管理员密码
export ADMIN_PASSWORD="your-strong-password"

# 3. 启动（生产模式）
MEMO_ENV=production \
GIN_MODE=release \
MEMO_JWT_SECRET="$JWT_SECRET" \
MEMO_ADMIN_PASSWORD="$ADMIN_PASSWORD" \
docker compose up -d --build
```

### 验证步骤

1. **检查健康端点：**
   ```bash
   curl http://localhost:9000/health
   ```

2. **检查日志：**
   ```bash
   docker compose logs -f
   ```

3. **检查 API 版本：**
   ```bash
   # 应该返回 v1
   curl http://localhost:9000/health
   ```

## 🐛 常见问题

### 端口被占用
```bash
lsof -i :9000
kill -9 <PID>
```

### 构建失败 - CGO 错误
确保 Docker 有足够的内存（至少 2GB）用于 Go 构建。

### 镜像过大
考虑使用 `docker buildx` 多架构构建，并启用构建缓存：
```bash
docker buildx build --platform linux/amd64,linux/arm64 -t memo-studio:latest --push .
```

## 📦 生产环境推荐配置

```yaml
# docker-compose.prod.yml
services:
  memo-studio:
    image: ghcr.io/yourusername/memo-studio:latest
    ports:
      - "9000:9000"
    environment:
      - MEMO_JWT_SECRET=${JWT_SECRET}
      - MEMO_ADMIN_PASSWORD=${ADMIN_PASSWORD}
      - MEMO_ENV=production
      - GIN_MODE=release
      - MEMO_CORS_ORIGINS=https://your.domain.com
    volumes:
      - memo_data:/data
    restart: unless-stopped
    deploy:
      resources:
        limits:
          cpus: '1'
          memory: 1G
```

## ✅ 部署后检查

- [ ] `curl http://localhost:9000/health` 返回 `{"status":"ok"}`
- [ ] JWT Secret 已设置为随机字符串
- [ ] 数据库文件在 `/data/notes.db`
- [ ] 附件目录在 `/data/storage`
- [ ] 日志无报错
- [ ] 前端可正常注册/登录
