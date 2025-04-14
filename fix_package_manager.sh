#!/bin/bash

# 检查是否为root用户
if [ "$EUID" -ne 0 ]; then
    echo "请使用root用户运行此脚本"
    exit 1
fi

# 设置镜像站点
MIRROR="https://dl.rockylinux.org/pub/rocky/9/BaseOS/x86_64/os/Packages"
TEMP_DIR="/tmp/pkg_fix"

# 创建临时目录
mkdir -p $TEMP_DIR
cd $TEMP_DIR

# 下载必要的包
echo "下载必要的包..."
wget $MIRROR/l/libsolv-0.7.22-1.el9.x86_64.rpm
wget $MIRROR/l/libdnf-0.69.0-1.el9.x86_64.rpm
wget $MIRROR/d/dnf-4.14.0-4.el9.noarch.rpm

# 安装包
echo "安装包..."
rpm -i --nodeps libsolv-0.7.22-1.el9.x86_64.rpm
rpm -i --nodeps libdnf-0.69.0-1.el9.x86_64.rpm
rpm -i --nodeps dnf-4.14.0-4.el9.noarch.rpm

# 清理
cd -
rm -rf $TEMP_DIR

echo "包管理器修复完成！" 