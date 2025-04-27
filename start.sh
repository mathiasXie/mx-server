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
        echo -e "\033[1;31mWarning: vosk_lib directory not found!\033[0m"
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

    # 检查Vad类库onnx是否存在，不存在则下载
    if [ ! -d "cgo_libs/onnxruntime_lib" ]; then
        echo -e "\033[1;31mWarning: onnxruntime directory not found!\033[0m"
        echo -e "\033[1;31mDownloading onnxruntime...\033[0m"
        # 根据CPU架构下载对应的Vosk类库
        arch=$(uname -m)
        if [[ "$arch" == "aarch64" ]]; then
            wget -O cgo_libs/onnxruntime-lib.tgz https://github.com/microsoft/onnxruntime/releases/download/v1.18.1/onnxruntime-linux-aarch64-1.18.1.tgz 
        else
            wget -O cgo_libs/onnxruntime-lib.tgz https://github.com/microsoft/onnxruntime/releases/download/v1.18.1/onnxruntime-linux-x64-1.18.1.tgz 
        fi

        tar -xzvf cgo_libs/onnxruntime-lib.tgz -C cgo_libs/
        mv cgo_libs/onnxruntime-linux* cgo_libs/onnxruntime_lib
        rm cgo_libs/onnxruntime-lib.tgz
    fi

    # 检查models目录是否存在，不存在则下载
    if [ ! -d "models/vosk-model-cn-0.22" ]; then
        echo -e "\033[1;31mWarning: vosk model directory not found!\033[0m"
        echo -e "\033[1;31mDownloading vosk model...\033[0m"
        wget -O models/vosk-model-cn-0.22.zip https://alphacephei.com/vosk/models/vosk-model-cn-0.22.zip
        unzip models/vosk-model-cn-0.22.zip -d models
        rm models/vosk-model-cn-0.22.zip
    fi

        # 检查models目录是否存在，不存在则下载
    #if [ ! -d "models/silero_vad.onnx" ]; then
    if [ ! -e "models/silero_vad.onnx" ]; then

        echo -e "\033[1;31mWarning: silero_vad.onnx not found!\033[0m"
        echo -e "\033[1;31mDownloading silero_vad.onnx...\033[0m"
        wget -O models/silero_vad.onnx https://github.com/snakers4/silero-vad/raw/master/src/silero_vad/data/silero_vad.onnx
    fi

    export VOSK_PATH="$(pwd)/cgo_libs/vosk_lib/"
    export ONNX_PATH="$(pwd)/cgo_libs/onnxruntime_lib"
    export LD_LIBRARY_PATH="$LD_LIBRARY_PATH:$VOSK_PATH:$ONNX_PATH/lib"
    export CGO_CPPFLAGS="-I $VOSK_PATH"
    export CGO_LDFLAGS="$CGO_LDFLAGS -L $VOSK_PATH"
    export CGO_LDFLAGS="$CGO_LDFLAGS -L $ONNX_PATH/lib -lonnxruntime"
    export CGO_CFLAGS="-I$ONNX_PATH/include"
    go run applications/asr-rpc/main.go -f applications/asr-rpc/conf &
    PIDS+=($!)
}

function START_TTS {
    echo -e "\033[1;33mStarting TTS Service...\033[0m"
    export SPEECHSDK_ROOT="$(pwd)/cgo_libs/microsoft_speechsdk" 

    # 检查cgo_libs/microsoft_speechsdk目录是否存在，不存在则下载
    if [ ! -d "$SPEECHSDK_ROOT" ]; then
        echo -e "\033[1;31mWarning: microsoft speechsdk directory not found!\033[0m"
        echo -e "\033[1;31mDownloading microsoft speechsdk...\033[0m"
        wget -O cgo_libs/SpeechSDK-Linux.tar.gz https://aka.ms/csspeech/linuxbinary
        tar -xzvf cgo_libs/SpeechSDK-Linux.tar.gz -C cgo_libs/
        rm cgo_libs/SpeechSDK-Linux.tar.gz
        mv cgo_libs/SpeechSDK-Linux* $SPEECHSDK_ROOT
    fi

    export CGO_CFLAGS="-I$SPEECHSDK_ROOT/include/c_api"
    export CGO_LDFLAGS="-L$SPEECHSDK_ROOT/lib/x64 -lMicrosoft.CognitiveServices.Speech.core"
    export LD_LIBRARY_PATH="$LD_LIBRARY_PATH:$SPEECHSDK_ROOT/lib/x64"
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

    # 需要安装 opus opus-devel opusfile opusfile-devel ogg 
    # 检查是否安装opus组件
    if pkg-config --cflags --libs opus > /dev/null 2>&1; then
        echo -e "\033[0;32mopus library installed successfully\033[0m"
        pkg-config --cflags --libs opus
    else
        echo -e "\033[0;31mError:opus library installation may have problems,please check the error information\033[0m"
        exit 1
    fi

    # 检查 gcc 是否安装
    if command -v gcc &> /dev/null; then
        echo -e "\033[0;32mgcc installed successfully\033[0m"
        gcc --version
    else
        echo -e "\033[0;31mError:gcc installation may have problems,please check the error information\033[0m"
        exit 1
    fi

    # 检查是否安装ffmpeg组件
    if command -v ffmpeg &> /dev/null; then
        echo -e "\033[0;32mffmpeg installed successfully\033[0m"
    else
        echo -e "\033[0;31mError:ffmpeg installation may have problems,please check the error information\033[0m"
        exit 1
    fi
    # 设置环境变量
    export PKG_CONFIG_PATH=$PKG_CONFIG_PATH:/usr/lib64/pkgconfig:/usr/local/lib/pkgconfig
    export LD_LIBRARY_PATH="/usr/lib64:/usr/lib:$LD_LIBRARY_PATH"
    echo "LD_LIBRARY_PATH:$LD_LIBRARY_PATH"
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
        START_ASR
        START_TTS
        START_FUNCTION
        START_LLM
        START_XIAOZHI
        ;;
esac

# 等待所有进程
wait
