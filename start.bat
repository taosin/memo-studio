@echo off
chcp 65001 >nul
setlocal enabledelayedexpansion

echo 🚀 启动 Memo Studio...

REM 检查 Go 是否安装
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 错误: 未找到 Go，请先安装 Go 1.21+
    echo    安装地址: https://go.dev/dl/
    pause
    exit /b 1
)

REM 检查 Node.js 是否安装
where node >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 错误: 未找到 Node.js，请先安装 Node.js
    echo    安装地址: https://nodejs.org/
    pause
    exit /b 1
)

REM 检查 npm 是否安装
where npm >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 错误: 未找到 npm，请先安装 npm
    pause
    exit /b 1
)

REM 获取脚本所在目录
cd /d "%~dp0"

REM 启动后端
echo 📦 启动后端服务...
cd backend

REM 安装 Go 依赖
if not exist "go.sum" (
    echo 📥 安装 Go 依赖...
    go mod download
)

REM 启动后端（后台运行）
start /B go run main.go > ..\backend.log 2>&1
cd ..

REM 等待后端启动
echo ⏳ 等待后端服务启动...
timeout /t 3 /nobreak >nul

REM 启动前端
echo 🎨 启动前端应用...
cd frontend

REM 安装前端依赖
if not exist "node_modules" (
    echo 📥 安装前端依赖（这可能需要几分钟）...
    call npm install
)

REM 启动前端
echo ✅ 服务已启动！
echo.
echo 📝 后端: http://localhost:9000
echo 🌐 前端: http://localhost:9001
echo.
echo 💡 提示: 关闭此窗口将停止所有服务
echo.

call npm run dev

cd ..
pause
