#!/bin/bash

# 检查是否为root用户
if [ "$EUID" -ne 0 ]; then
    echo "请使用root用户运行此脚本"
    exit 1
fi

# 设置临时目录
TEMP_DIR="/tmp/lsof_install"
mkdir -p $TEMP_DIR
cd $TEMP_DIR

# 下载lsof包
echo "下载lsof包..."
wget https://dl.rockylinux.org/pub/rocky/9/BaseOS/x86_64/os/Packages/l/lsof-4.94.0-1.el9.x86_64.rpm

# 安装lsof
echo "安装lsof..."
rpm -i --nodeps lsof-4.94.0-1.el9.x86_64.rpm

# 验证安装
if command -v lsof >/dev/null 2>&1; then
    echo "lsof安装成功！"
    lsof --version
else
    echo "lsof安装失败，请检查错误信息"
    exit 1
fi

# 清理
cd -
rm -rf $TEMP_DIR

echo "安装完成！" 