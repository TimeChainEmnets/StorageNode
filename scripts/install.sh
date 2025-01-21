#!/bin/bash

# 检查Go环境
if ! command -v go &> /dev/null; then
    echo "未找到Go环境，请先安装Go"
    exit 1
fi

# 检查IPFS环境
if ! command -v ipfs &> /dev/null; then
    echo "未找到IPFS，正在安装..."
    # 这里可以添加安装IPFS的命令
fi

# 设置默认仓库路径
DEFAULT_REPO_PATH="$HOME/.myipfs"

# 创建仓库目录
if [ ! -d "$DEFAULT_REPO_PATH" ]; then
    mkdir -p "$DEFAULT_REPO_PATH"
    echo "创建仓库目录: $DEFAULT_REPO_PATH"
fi

# 编译和安装
echo "正在编译项目..."
go build -o myipfs ./cmd/myipfs

# 移动可执行文件到系统路径
sudo mv myipfs /usr/local/bin/

echo "安装完成！"
echo "你可以通过运行 'myipfs' 命令来启动程序"