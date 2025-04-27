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

echo -e "${GREEN}>>> 检测系统类型...${RESET}"

OS_TYPE="$(uname -s)"

case "$OS_TYPE" in
    Linux)
        if [ -f /etc/os-release ]; then
            . /etc/os-release
            OS=$ID
            VERSION_ID=$VERSION_ID
        else
            echo "无法识别 Linux 系统，缺少 /etc/os-release"
            exit 1
        fi
        ;;
    Darwin)
        OS="macos"
        ;;
    *)
        echo -e "${YELLOW}不支持的系统类型: $OS_TYPE${RESET}"
        exit 1
        ;;
esac

# === Red Hat 系系统 ===
if [[ "$OS" =~ (rhel|centos|rocky|almalinux|fedora) ]]; then
    echo -e "${GREEN}>>> 检测到 RedHat 系系统：$OS${RESET}"
    PM=$(command -v dnf || command -v yum)

    echo -e "${GREEN}>>> 安装 EPEL 和 RPMFusion 源${RESET}"
    $PM install -y epel-release || true
    MAJOR_VERSION=${VERSION_ID%%.*}
    $PM install -y \
        "https://mirrors.rpmfusion.org/free/el/rpmfusion-free-release-${MAJOR_VERSION}.noarch.rpm" \
        "https://mirrors.rpmfusion.org/nonfree/el/rpmfusion-nonfree-release-${MAJOR_VERSION}.noarch.rpm" || true

    echo -e "${GREEN}>>> 更新缓存并安装依赖（opus, opusfile, ogg, gcc, ffmpeg, zip）${RESET}"
    $PM makecache -y
    $PM install -y opus opus-devel opusfile opusfile-devel libogg-devel gcc ffmpeg ffmpeg-devel pkgconfig zip unzip

# === Debian 系系统 ===
elif [[ "$OS" =~ (debian|ubuntu|linuxmint) ]]; then
    echo -e "${GREEN}>>> 检测到 Debian 系系统：$OS${RESET}"
    apt update -y
    apt install -y libopus-dev libopusfile-dev libogg-dev gcc ffmpeg pkg-config zip unzip

# === macOS 系统 ===
elif [[ "$OS" == "macos" ]]; then
    echo -e "${GREEN}>>> 检测到 macOS 系统${RESET}"
    if ! command -v brew >/dev/null 2>&1; then
        echo -e "${YELLOW}>>> 未检测到 Homebrew，开始安装...${RESET}"
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        eval "$(/opt/homebrew/bin/brew shellenv)"  # 支持 M1/M2
    fi

    echo -e "${GREEN}>>> 使用 Homebrew 安装 opus、opusfile、ogg、gcc、ffmpeg、zip${RESET}"
    brew install opus opusfile libogg gcc ffmpeg pkg-config zip unzip
else
    echo -e "${YELLOW}未知或未适配的系统: $OS${RESET}"
    exit 1
fi


# === 通用验证 ===
echo -e "${YELLOW}>>> 验证 opusfile 是否可用${RESET}"
if ! pkg-config --modversion opusfile >/dev/null 2>&1; then
    echo -e "${YELLOW}警告：未找到 opusfile.pc，请检查 PKG_CONFIG_PATH 设置${RESET}"
    echo -e "建议添加：${CYAN}export PKG_CONFIG_PATH=/usr/lib64/pkgconfig:/usr/lib/pkgconfig:/usr/local/lib/pkgconfig:/opt/homebrew/lib/pkgconfig${RESET}"
else
    echo -e "${GREEN}pkg-config 成功识别 opusfile: $(pkg-config --modversion opusfile)${RESET}"
fi

# === 验证 ffmpeg 是否可用 === ffmpeg不是pkg包，是单独的命令
echo -e "${YELLOW}>>> 验证 ffmpeg 是否可用${RESET}"
if ! command -v ffmpeg -version >/dev/null 2>&1; then
    echo -e "${YELLOW}警告：未找到 ffmpeg，请检查是否安装${RESET}"
else
    echo -e "${GREEN}ffmpeg 成功识别${RESET}"
fi

echo -e "${YELLOW}当前 pkg-config 搜索路径：$(pkg-config --variable pc_path pkg-config)${RESET}"
echo -e "${GREEN}>>> 安装完成 ✅${RESET}"
