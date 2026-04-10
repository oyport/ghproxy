@echo off
REM GHProxy 管理后台快速启动脚本 (Windows)

echo =========================================
echo GHProxy 管理后台启动脚本
echo =========================================

REM 检查配置文件
if not exist "config\config.toml" (
    echo 错误: 配置文件不存在
    pause
    exit /b 1
)

REM 编译项目
echo 正在编译项目...
go build -o ghproxy.exe main.go

if %errorlevel% neq 0 (
    echo 编译失败
    pause
    exit /b 1
)

echo 编译成功

REM 启动服务
echo 正在启动服务...
ghproxy.exe -c config\config.toml

pause
