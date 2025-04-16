#!/bin/bash

# 存储所有后台进程的PID
declare -a PIDS

# 定义清理函数
cleanup() {
    echo "Stopping all processes..."
    for pid in "${PIDS[@]}"; do
        if kill -0 $pid 2>/dev/null; then
            kill $pid
        fi
    done
    exit 0
}

# 设置信号处理
trap cleanup SIGINT SIGTERM

# 设置环境变量

#微软Azure tts使用环境变量
export SPEECHSDK_ROOT="$HOME/speechsdk" 
export CGO_CFLAGS="-I$SPEECHSDK_ROOT/include/c_api"
export CGO_LDFLAGS="-L$SPEECHSDK_ROOT/lib/x64 -lMicrosoft.CognitiveServices.Speech.core"
export LD_LIBRARY_PATH="$SPEECHSDK_ROOT/lib/x64:$LD_LIBRARY_PATH"

export LD_LIBRARY_PATH="/usr/lib64:/usr/lib:$LD_LIBRARY_PATH"

# export VOSK_PATH="./applications/asr-rpc/vosk_lib/"

# export LD_LIBRARY_PATH="/usr/lib64:/usr/lib:$LD_LIBRARY_PATH"

# 确保系统库路径在LD_LIBRARY_PATH中

# 启动所有服务
go run applications/asr-rpc/main.go -f applications/asr-rpc/conf &
PIDS+=($!)
go run applications/tts-rpc/main.go -f applications/tts-rpc/conf &
PIDS+=($!)
go run applications/function-rpc/main.go -f applications/function-rpc/conf &
PIDS+=($!)
go run applications/llm-rpc/main.go -f applications/llm-rpc/conf &
PIDS+=($!)
go run applications/xiaozhi-server/main.go -f applications/xiaozhi-server/conf &
PIDS+=($!)

# 等待所有进程
wait

