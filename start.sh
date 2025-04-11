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

# 启动所有服务
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

