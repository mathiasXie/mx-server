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
export LD_LIBRARY_PATH="/usr/lib64:/usr/lib:$LD_LIBRARY_PATH"


if [ ! -d "models" ]; then
    mkdir models
    echo "*" > models/.gitignore
fi
if [ ! -d "cgo_libs" ]; then
    mkdir cgo_libs
    echo "*" > cgo_libs/.gitignore
fi


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


function START_ASR {
    echo -e "\033[1;31mStarting ASR Service...\033[0m"

    # 检查Vosk类库是否存在，不存在则下载
    if [ ! -d "cgo_libs/vosk_lib" ]; then
        echo -e "\033[1;31mError: vosk_lib directory not found!\033[0m"
        echo -e "\033[1;31mDownloading vosk_lib...\033[0m"
        # 根据CPU架构下载对应的Vosk类库
        arch=$(uname -m)
        if [[ "$arch" == "aarch64" ]]; then
            wget -O cgo_libs/vosk-lib.zip https://github.com/alphacep/vosk-api/releases/download/v0.3.45/vosk-linux-aarch64-0.3.45.zip 
        else
            wget -O cgo_libs/vosk-lib.zip https://github.com/alphacep/vosk-api/releases/download/v0.3.45/vosk-linux-x86_64-0.3.45.zip 
        fi
        
        unzip cgo_libs/vosk-lib.zip -d cgo_libs/
        mv cgo_libs/vosk-linux* cgo_libs/vosk_lib
        rm cgo_libs/vosk-lib.zip
    fi

    # 检查models目录是否存在，不存在则下载
    if [ ! -d "models/vosk-model-cn-0.22" ]; then
        echo -e "\033[1;31mError: vosk model directory not found!\033[0m"
        echo -e "\033[1;31mDownloading vosk model...\033[0m"
        wget -O models/vosk-model-cn-0.22.zip https://alphacephei.com/vosk/models/vosk-model-cn-0.22.zip
        unzip models/vosk-model-cn-0.22.zip -d models
        rm models/vosk-model-cn-0.22.zip
    fi

    export VOSK_PATH="$(pwd)/cgo_libs/vosk_lib/"
    export LD_LIBRARY_PATH="$VOSK_PATH:$LD_LIBRARY_PATH"
    export CGO_CPPFLAGS="-I $VOSK_PATH"
    export CGO_LDFLAGS="$CGO_LDFLAGS -L $VOSK_PATH"
    go run applications/asr-rpc/main.go -f applications/asr-rpc/conf &
    PIDS+=($!)
}

function START_TTS {
    echo -e "\033[1;33mStarting TTS Service...\033[0m"
    export SPEECHSDK_ROOT="$(pwd)/cgo_libs/microsoft_speechsdk" 

    # 检查cgo_libs/microsoft_speechsdk目录是否存在，不存在则下载
    if [ ! -d "$SPEECHSDK_ROOT" ]; then
        echo -e "\033[1;31mError: microsoft speechsdk directory not found!\033[0m"
        echo -e "\033[1;31mDownloading microsoft speechsdk...\033[0m"
        wget -O cgo_libs/SpeechSDK-Linux.tar.gz https://aka.ms/csspeech/linuxbinary
        tar -xzvf cgo_libs/SpeechSDK-Linux.tar.gz -C cgo_libs/
        rm cgo_libs/SpeechSDK-Linux.tar.gz
        mv cgo_libs/SpeechSDK-Linux* $SPEECHSDK_ROOT
    fi

    export CGO_CFLAGS="-I$SPEECHSDK_ROOT/include/c_api"
    export CGO_LDFLAGS="-L$SPEECHSDK_ROOT/lib/x64 -lMicrosoft.CognitiveServices.Speech.core"
    export LD_LIBRARY_PATH="$SPEECHSDK_ROOT/lib/x64:$LD_LIBRARY_PATH"
    go run applications/tts-rpc/main.go -f applications/tts-rpc/conf &
    PIDS+=($!)
}

function START_FUNCTION {
    echo -e "\033[1;35mStarting Function Service...\033[0m"
    go run applications/function-rpc/main.go -f applications/function-rpc/conf &
    PIDS+=($!)
}

function START_LLM {
    echo -e "\033[1;36mStarting LLM Service...\033[0m"
    go run applications/llm-rpc/main.go -f applications/llm-rpc/conf &
    PIDS+=($!)
}

function START_XIAOZHI {
    echo -e "\033[1;37mStarting XiaoZhi Server...\033[0m"

    # 检查是否安装 ffmpeg
    if ! command -v ffmpeg &> /dev/null; then
        echo -e "\033[1;31mError: ffmpeg could not be found!\033[0m"
        echo -e "\033[1;31mPlease install ffmpeg first.\033[0m"
        exit 1
    fi
    go run applications/xiaozhi-server/main.go -f applications/xiaozhi-server/conf &
    PIDS+=($!)
}

# 启动服务
case $choice in
    1)
        START_ASR
        ;;
    2)
        START_TTS
        ;;
    3)
        START_FUNCTION
        ;;
    4)
        START_LLM
        ;;
    5)
        START_XIAOZHI
        ;;
    6)
        echo -e "\033[1;32mStarting All Services...\033[0m"
        for i in {1..5}; do
            case $i in
                1) START_ASR ;;
                2) START_TTS ;;
                3) START_FUNCTION ;;
                4) START_LLM ;;
                5) START_XIAOZHI ;;
            esac
        done
        ;;
esac

# 等待所有进程
wait
