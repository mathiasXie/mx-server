
#!/bin/bash

set -e

GREEN="\033[1;32m"
YELLOW="\033[1;33m"
CYAN="\033[0;36m"
RESET="\033[0m"

# 检查是否为 root 用户
if [ "$EUID" -ne 0 ]; then
  echo -e "${YELLOW}请使用 root 或 sudo 权限运行此脚本${RESET}"
  exit 1
fi

# 检查是否已安装 go
if command -v go >/dev/null 2>&1; then
    echo -e "${GREEN}go 已安装${RESET}"
    exit 0
fi

# 下载 go1.24.1 安装包
echo -e "${GREEN}>>> 下载 go1.24.1 安装包${RESET}"
wget https://go.dev/dl/go1.24.1.linux-amd64.tar.gz

# 解压安装包
echo -e "${GREEN}>>> 解压安装包${RESET}"
tar -C /usr/local -xzf go1.24.1.linux-amd64.tar.gz

# 设置环境变量
echo -e "${GREEN}>>> 设置环境变量${RESET}"
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:/usr/local/go/bin:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc



# 验证安装
echo -e "${GREEN}>>> 验证安装${RESET}"
go version

echo -e "${GREEN}>>> 安装完成${RESET}"

