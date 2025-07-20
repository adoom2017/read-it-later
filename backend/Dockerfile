# 使用官方的 Go 镜像作为构建环境
FROM golang:1.24.4-alpine AS builder

# 安装构建依赖
RUN apk add --no-cache gcc musl-dev sqlite-dev

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main . && \
    ls -la /app/main

# 使用 Alpine 镜像作为运行环境（与构建环境保持一致）
FROM alpine:latest

# 安装运行时依赖和 Chrome
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    sqlite \
    curl \
    chromium \
    chromium-chromedriver \
    xvfb-run \
    font-noto-emoji \
    ttf-freefont \
    && rm -rf /var/cache/apk/*

# 设置时区环境变量
ENV TZ=UTC

# 创建非root用户
RUN addgroup -g 1001 appgroup && \
    adduser -D -u 1001 -G appgroup -s /bin/sh appuser

# 设置工作目录
WORKDIR /app

# 从构建阶段复制可执行文件
COPY --from=builder /app/main .

# 创建数据目录并设置权限
RUN mkdir -p /app/data && \
    chown -R appuser:appgroup /app

# 切换到非root用户
USER appuser

# 设置环境变量
ENV GIN_MODE=release
ENV PORT=8080

# 暴露端口
EXPOSE 8080

# 健康检查
HEALTHCHECK --interval=30s --timeout=10s --start-period=30s --retries=3 \
  CMD curl -f http://localhost:8080/ || exit 1

# 运行应用
CMD ["./main"]
