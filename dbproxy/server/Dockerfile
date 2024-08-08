# 使用官方 Go 镜像作为构建阶段的基础镜像
FROM golang:1.20 AS builder

# 设置工作目录
WORKDIR /app

# 复制 go.mod 和 go.sum 文件
COPY ../go.mod ../go.sum ./
# 下载项目的依赖
RUN go mod tidy

# 复制源代码
COPY . .

# 编译 Go 应用
RUN CGO_ENABLED=0 GOOS=linux go build -o dbproxy cmd/main.go

# 使用一个轻量级的镜像作为运行阶段的基础镜像
FROM alpine:latest

# 安装必要的包（如果你的应用需要特定的系统包，可以在这里添加）
RUN apk --no-cache add ca-certificates

# 设置工作目录
WORKDIR /root/

# 从构建阶段复制编译好的二进制文件
COPY --from=builder /app/dbproxy .

# 暴露端口
EXPOSE 10000

# 设置启动命令
CMD ["./dbproxy"]