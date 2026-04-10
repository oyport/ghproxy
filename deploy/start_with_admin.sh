#!/bin/bash

# GHProxy 管理后台快速启动脚本

echo "========================================="
echo "GHProxy 管理后台启动脚本"
echo "========================================="

# 检查配置文件
if [ ! -f "config/config.toml" ]; then
    echo "错误: 配置文件不存在"
    exit 1
fi

# 检查是否已启用管理后台
if grep -q "enabled = true" config/config.toml | grep -A2 "\[admin\]"; then
    echo "管理后台已启用"
else
    echo "警告: 管理后台未启用"
    echo "请修改 config/config.toml 文件，设置 [admin] enabled = true"
    read -p "是否现在启用管理后台? (y/n): " choice
    if [ "$choice" = "y" ]; then
        sed -i 's/^\(\[admin\]\)/\1\nenabled = true/' config/config.toml
        echo "管理后台已启用"
    fi
fi

# 编译项目
echo "正在编译项目..."
go build -o ghproxy main.go

if [ $? -ne 0 ]; then
    echo "编译失败"
    exit 1
fi

echo "编译成功"

# 启动服务
echo "正在启动服务..."
./ghproxy -c config/config.toml
