#!/bin/bash

# 检查是否为root用户
if [ "$EUID" -ne 0 ]; then
    echo "请使用root用户运行此脚本"
    exit 1
fi

# 检查系统类型
if [ -f /etc/os-release ]; then
    . /etc/os-release
    OS=$NAME
    VER=$VERSION_ID
else
    echo "无法确定操作系统类型"
    exit 1
fi

echo "检测到系统: $OS $VER"

# 尝试修复包管理器
echo "尝试修复包管理器..."
if command -v microdnf >/dev/null 2>&1; then
    echo "使用microdnf包管理器..."
    PKG_MANAGER="microdnf"
elif command -v yum >/dev/null 2>&1; then
    echo "使用yum包管理器..."
    PKG_MANAGER="yum"
elif command -v dnf >/dev/null 2>&1; then
    echo "使用dnf包管理器..."
    PKG_MANAGER="dnf"
else
    echo "未找到可用的包管理器，尝试安装microdnf..."
    rpm -i --nodeps https://dl.rockylinux.org/pub/rocky/9/BaseOS/x86_64/os/Packages/m/microdnf-3.8.0-2.el9.noarch.rpm
    PKG_MANAGER="microdnf"
fi

# 安装EPEL仓库
echo "正在安装EPEL仓库..."
$PKG_MANAGER install -y epel-release || {
    echo "安装EPEL仓库失败，尝试直接下载..."
    wget https://dl.fedoraproject.org/pub/epel/epel-release-latest-9.noarch.rpm
    rpm -i --nodeps epel-release-latest-9.noarch.rpm
    rm -f epel-release-latest-9.noarch.rpm
}

# 启用CRB仓库
echo "正在启用CRB仓库..."
if [ -f /usr/bin/crb ]; then
    /usr/bin/crb enable || {
        echo "启用CRB仓库失败，尝试使用替代方法..."
        $PKG_MANAGER config-manager --set-enabled crb || {
            echo "启用CRB仓库失败"
            exit 1
        }
    }
else
    echo "未找到crb命令，尝试使用替代方法..."
    $PKG_MANAGER config-manager --set-enabled crb || {
        echo "启用CRB仓库失败"
        exit 1
    }
fi

# 安装Opus相关依赖
echo "正在安装Opus相关依赖..."
$PKG_MANAGER install -y opus-devel opusfile-devel || {
    echo "安装Opus相关依赖失败，尝试直接下载..."
    wget https://dl.rockylinux.org/pub/rocky/9/CRB/x86_64/os/Packages/o/opus-devel-1.3.1-10.el9.x86_64.rpm
    wget https://dl.rockylinux.org/pub/rocky/9/CRB/x86_64/os/Packages/o/opusfile-devel-0.12-6.el9.x86_64.rpm
    rpm -i --nodeps opus-devel-1.3.1-10.el9.x86_64.rpm opusfile-devel-0.12-6.el9.x86_64.rpm
    rm -f opus-devel-1.3.1-10.el9.x86_64.rpm opusfile-devel-0.12-6.el9.x86_64.rpm
}

# 验证安装
echo "验证安装..."
if pkg-config --cflags --libs opus > /dev/null 2>&1; then
    echo "Opus库安装成功！"
    pkg-config --cflags --libs opus
else
    echo "Opus库安装可能存在问题，请检查错误信息"
    exit 1
fi

echo "安装完成！" 