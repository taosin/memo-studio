# 后端环境设置

## 解决依赖问题

如果遇到 `missing go.sum entry` 错误，请按以下步骤操作：

### 方法一：自动安装（推荐）

在项目根目录运行启动脚本，脚本会自动安装依赖：

```bash
./start.sh
```

### 方法二：手动安装

如果自动安装失败，请手动执行：

```bash
cd backend
go mod download
go mod tidy
```

### 方法三：使用 Go 代理（推荐，解决网络问题）

如果遇到网络超时或无法访问默认的 Go 代理，请设置国内镜像：

```bash
# 方法1：临时设置（当前终端有效）
export GOPROXY=https://goproxy.cn,direct

# 方法2：使用阿里云镜像（如果 goproxy.cn 不可用）
export GOPROXY=https://mirrors.aliyun.com/goproxy/,direct

# 然后安装依赖
cd backend
go mod download
go mod tidy
```

**永久设置（推荐）：**

```bash
# 添加到 ~/.zshrc（如果使用 zsh）
echo 'export GOPROXY=https://goproxy.cn,direct' >> ~/.zshrc
source ~/.zshrc

# 或添加到 ~/.bash_profile（如果使用 bash）
echo 'export GOPROXY=https://goproxy.cn,direct' >> ~/.bash_profile
source ~/.bash_profile
```

**注意：** 启动脚本 `start.sh` 已自动设置代理，无需手动配置。

### 验证安装

安装完成后，检查 `go.sum` 文件是否存在：

```bash
ls -la backend/go.sum
```

如果文件存在，说明依赖已正确安装。

## 常见问题

### 1. 网络连接问题

如果遇到 `dial tcp: lookup proxy.golang.org` 错误：

- 检查网络连接
- 设置 Go 代理（见方法三）
- 使用 VPN（如果需要）

### 2. 权限问题

确保有写入 `backend` 目录的权限：

```bash
chmod -R 755 backend
```

### 3. Go 版本问题

确保 Go 版本 >= 1.21：

```bash
go version
```

如果版本过低，请升级 Go：https://go.dev/dl/
