#!/bin/bash

# 存储所有后台进程的PID
declare -a PIDS

# 定义清理函数
cleanup() {
    echo -e "\n\033[1;31mStopping all processes...\033[0m"
    for pid in "${PIDS[@]}"; do
        if kill -0 $pid 2>/dev/null; then
            kill $pid
            echo -e "\033[1;33mStopped process with PID: $pid\033[0m"
        fi
    done
    exit 0
}

# 设置信号处理
trap cleanup SIGINT SIGTERM

# 设置环境变量
export SPEECHSDK_ROOT="$HOME/speechsdk" 
export CGO_CFLAGS="-I$SPEECHSDK_ROOT/include/c_api"
export CGO_LDFLAGS="-L$SPEECHSDK_ROOT/lib/x64 -lMicrosoft.CognitiveServices.Speech.core"
export LD_LIBRARY_PATH="$SPEECHSDK_ROOT/lib/x64:$LD_LIBRARY_PATH"

export LD_LIBRARY_PATH="/usr/lib64:/usr/lib:$LD_LIBRARY_PATH"

export VOSK_PATH="$(pwd)/applications/asr-rpc/vosk_lib/"
export LD_LIBRARY_PATH="$VOSK_PATH:$LD_LIBRARY_PATH"
export CGO_CPPFLAGS="-I $VOSK_PATH"
export CGO_LDFLAGS="$CGO_LDFLAGS -L $VOSK_PATH"

# 服务列表
services=(
    "\033[0;31m1) ASR Service\033[0m"
    "\033[0;33m2) TTS Service\033[0m"
    "\033[0;35m3) Function Service\033[0m"
    "\033[0;36m4) LLM Service\033[0m"
    "\033[1;33m5) XiaoZhi Server\033[0m"
    "\033[1;32m6) All Services\033[0m"  # 绿色
)

# 处理参数
while getopts "s:" opt; do
    case $opt in
        s)
            choice=$OPTARG
            ;;
        *)
            echo -e "\033[1;31mInvalid option! Exiting...\033[0m"
            exit 1
            ;;
    esac
done

# 显示服务选项
echo -e "\033[1;34mSelect a service to start:\033[0m"
for service in "${services[@]}"; do
    echo -e "$service"
done

# 如果没有通过参数指定选择，提示用户输入
if [[ $choice -lt 1 || $choice -gt 6 ]]; then
    read -p "Enter your choice (1-6, default is 6): " input_choice
    choice=${input_choice:-6}  # 如果用户没有输入，则使用默认值6
fi

# 启动服务
case $choice in
    1)
        echo -e "\033[1;32mStarting ASR Service...\033[0m"
        go run applications/asr-rpc/main.go -f applications/asr-rpc/conf &
        PIDS+=($!)
        ;;
    2)
        echo -e "\033[1;32mStarting TTS Service...\033[0m"
        go run applications/tts-rpc/main.go -f applications/tts-rpc/conf &
        PIDS+=($!)
        ;;
    3)
        echo -e "\033[1;32mStarting Function Service...\033[0m"
        go run applications/function-rpc/main.go -f applications/function-rpc/conf &
        PIDS+=($!)
        ;;
    4)
        echo -e "\033[1;32mStarting LLM Service...\033[0m"
        go run applications/llm-rpc/main.go -f applications/llm-rpc/conf &
        PIDS+=($!)
        ;;
    5)
        echo -e "\033[1;32mStarting XiaoZhi Server...\033[0m"
        go run applications/xiaozhi-server/main.go -f applications/xiaozhi-server/conf &
        PIDS+=($!)
        ;;
    6)
        echo -e "\033[1;32mStarting All Services...\033[0m"
        for i in {1..5}; do
            case $i in
                1) go run applications/asr-rpc/main.go -f applications/asr-rpc/conf & PIDS+=($!) ;;
                2) go run applications/tts-rpc/main.go -f applications/tts-rpc/conf & PIDS+=($!) ;;
                3) go run applications/function-rpc/main.go -f applications/function-rpc/conf & PIDS+=($!) ;;
                4) go run applications/llm-rpc/main.go -f applications/llm-rpc/conf & PIDS+=($!) ;;
                5) go run applications/xiaozhi-server/main.go -f applications/xiaozhi-server/conf & PIDS+=($!) ;;
            esac
        done
        ;;
esac

# 等待所有进程
wait
