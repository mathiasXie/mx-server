# 构建阶段
FROM golang:1.23-alpine AS builder

WORKDIR /app

# 设置 Go 环境变量
ENV GOPROXY='https://goproxy.cn,direct'

# 复制 go.mod 和 go.sum
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -o gin-web github.com/mathiasXie/gin-web/applications/gin-web

# 最终阶段
FROM alpine:3.21

# 安装基本工具
# RUN apk add --no-cache curl tar bash lib c6-compat

# 更新软件源索引
RUN apk update

# 安装所需的软件包，假设使用 libc6-compat 替代 c6-compat
RUN apk add --no-cache curl tar bash libc-dev libc6-compat
# 下载并安装预编译的 Filebeat
# https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-7.17.28-linux-arm64.tar.gz
ENV FILEBEAT_VERSION=7.17.28
ENV TZ Asia/Shanghai

RUN curl -L -O https://artifacts.elastic.co/downloads/beats/filebeat/filebeat-${FILEBEAT_VERSION}-linux-arm64.tar.gz && \
    tar xzvf filebeat-${FILEBEAT_VERSION}-linux-arm64.tar.gz && \
    mv filebeat-${FILEBEAT_VERSION}-linux-arm64/filebeat /usr/local/bin/ && \
    rm -rf filebeat-${FILEBEAT_VERSION}-linux-arm64 filebeat-${FILEBEAT_VERSION}-linux-arm64.tar.gz
ENV TZ Asia/Shanghai

# 创建必要的目录
RUN mkdir -p /app/logs /app/conf /var/log/filebeat

# 从构建阶段复制应用
COPY --from=builder /app/gin-web /app/gin-web
COPY filebeat.yml /etc/filebeat/filebeat.yml
COPY conf/gin_web_prod.yaml /app/conf/gin_web_prod.yaml

# 设置工作目录
WORKDIR /app

# 设置环境变量
ENV SERVICE_NAME=gin-web \
    ENV=production \
    LOGSTASH_HOST=logstash \
    LOGSTASH_PORT=5044

# 创建启动脚本
RUN echo '#!/bin/sh' > /app/start.sh && \
    echo 'chmod go-w /etc/filebeat/filebeat.yml' >> /app/start.sh && \
    echo 'filebeat -e -c /etc/filebeat/filebeat.yml &' >> /app/start.sh && \
    echo '/app/gin-web -f /app/conf -env prod' >> /app/start.sh && \
    #echo 'more /app/logs/server-latest.log' >> /app/start.sh && \
    #echo 'more /var/log/filebeat/filebeat.log' >> /app/start.sh && \
    #echo 'more /etc/filebeat/filebeat.yml' >> /app/start.sh && \

    chmod +x /app/start.sh && \
    sed -i 's/\r$//' /app/start.sh

# 暴露应用端口
EXPOSE 8080

# 启动应用和 Filebeat
CMD ["/app/start.sh"] 