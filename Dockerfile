FROM alpine:latest

# 替换软件源为腾讯云镜像源
RUN sed -i 's|https://dl-cdn.alpinelinux.org/alpine|http://mirrors.tencentyun.com/alpine|g' /etc/apk/repositories


# 安装时区工具
RUN apk update && \
    apk upgrade --no-cache && \
    apk add --no-cache tzdata && \
    rm -rf /var/cache/apk/*

# 设置上海时区
ENV TZ=Asia/Shanghai

ENV GF_GCFG_FILE=config.prod.yaml

# 复制 Go 应用程序到容器中
COPY hiveify-core /app/

# 设置可执行权限
RUN chmod +x /app/hiveify-core

# 暴露端口
EXPOSE 8080

# 设置工作目录
WORKDIR /app

# 运行 Go 应用程序
CMD ["./hiveify-core"]
