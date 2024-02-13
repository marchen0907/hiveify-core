FROM debian:latest

# 替换软件源为腾讯云镜像源
RUN echo "Types: deb" > /etc/apt/sources.list.d/debian.sources && \
    echo "URIs: http://mirrors.tencentyun.com/debian" >> /etc/apt/sources.list.d/debian.sources && \
    echo "Suites: bookworm" >> /etc/apt/sources.list.d/debian.sources && \
    echo "Components: main" >> /etc/apt/sources.list.d/debian.sources && \
    echo "Signed-By: /usr/share/keyrings/debian-archive-keyring.gpg" >> /etc/apt/sources.list.d/debian.sources

RUN echo "Types: deb" > /etc/apt/sources.list.d/security.sources && \
    echo "URIs: http://mirrors.tencentyun.com/debian-security" >> /etc/apt/sources.list.d/security.sources && \
    echo "Suites: bookworm-security" >> /etc/apt/sources.list.d/security.sources && \
    echo "Components: main" >> /etc/apt/sources.list.d/security.sources && \
    echo "Signed-By: /usr/share/keyrings/debian-archive-keyring.gpg" >> /etc/apt/sources.list.d/security.sources

# 安装基本工具和时区信息
RUN apt update && \
    apt install -y \
        tzdata \
        ca-certificates \
        && \
    rm -rf /var/lib/apt/lists/*

# 设置上海时区
ENV TZ=Asia/Shanghai

# 设置环境变量
ENV ENV_FILE .env.prod

# 复制 Go 应用程序和 .env.prod 到容器中
COPY hiveify-core .env.prod /app/

RUN chmod +x /app/hiveify-core

# 设置工作目录
WORKDIR /app

# 运行 Go 应用程序
CMD ["./hiveify-core"]
