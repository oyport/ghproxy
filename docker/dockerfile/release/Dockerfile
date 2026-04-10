# GHProxy Release Dockerfile
# 多阶段构建，优化镜像大小

# 构建阶段
FROM alpine:latest AS builder

# 构建参数
ARG USER=WJQSERVER-STUDIO
ARG REPO=ghproxy
ARG BRANCH=main
ARG APPLICATION=ghproxy
ARG TARGETOS=linux
ARG TARGETARCH=amd64

# 创建必要的目录
RUN mkdir -p /data/www && \
    mkdir -p /data/${APPLICATION}/config && \
    mkdir -p /data/${APPLICATION}/log

# 安装依赖
RUN apk add --no-cache curl wget tar

# 下载并解压应用程序
RUN VERSION=$(curl -s https://raw.githubusercontent.com/${USER}/${REPO}/${BRANCH}/VERSION) && \
    echo "Downloading ${APPLICATION} version: ${VERSION}" && \
    wget -O /data/${APPLICATION}/${APPLICATION}-${TARGETOS}-${TARGETARCH}.tar.gz \
        https://github.com/${USER}/${REPO}/releases/download/${VERSION}/${APPLICATION}-${TARGETOS}-${TARGETARCH}.tar.gz && \
    tar -zxvf /data/${APPLICATION}/${APPLICATION}-${TARGETOS}-${TARGETARCH}.tar.gz -C /data/${APPLICATION} && \
    rm -rf /data/${APPLICATION}/${APPLICATION}-${TARGETOS}-${TARGETARCH}.tar.gz

# 下载配置文件
RUN wget -O /data/${APPLICATION}/config/config.toml \
        https://raw.githubusercontent.com/${USER}/${REPO}/${BRANCH}/config/config.toml && \
    wget -O /data/${APPLICATION}/config/blacklist.json \
        https://raw.githubusercontent.com/${USER}/${REPO}/${BRANCH}/config/blacklist.json && \
    wget -O /data/${APPLICATION}/config/whitelist.json \
        https://raw.githubusercontent.com/${USER}/${REPO}/${BRANCH}/config/whitelist.json

# 设置权限
RUN chmod +x /data/${APPLICATION}/${APPLICATION}

# 运行阶段
FROM alpine:latest

# 安装运行时依赖
RUN apk add --no-cache curl tzdata

# 构建参数
ARG APPLICATION=ghproxy

# 设置工作目录
WORKDIR /data/${APPLICATION}

# 从构建阶段复制文件
COPY --from=builder /data/www /data/www
COPY --from=builder /data/${APPLICATION} /data/${APPLICATION}

# 创建必要的目录并设置权限
RUN chmod +x /data/${APPLICATION}/${APPLICATION} && \
    mkdir -p /data/${APPLICATION}/log && \
    mkdir -p /data/${APPLICATION}/config

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8080/api/healthcheck || exit 1

# 设置环境变量
ENV TZ=Asia/Shanghai

# 运行应用
CMD ["/data/ghproxy/ghproxy", "-c", "/data/ghproxy/config/config.toml"]
