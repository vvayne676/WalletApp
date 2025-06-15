# 💰 WalletApp 简化版 Dockerfile

FROM public.ecr.aws/docker/library/golang:1.24-alpine3.21

# 设置工作目录
WORKDIR /app

# 复制所有文件
COPY . .

# 下载依赖并构建
RUN go build -o walletapp

# 暴露端口
EXPOSE 8080

# 启动应用
CMD ["./walletapp"] 