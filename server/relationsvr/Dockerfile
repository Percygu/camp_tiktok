FROM zhugeqing/golang:1.19  AS builder
WORKDIR /app
MAINTAINER 诸葛青
COPY . .
# 设置代理环境变量
ARG GOPROXY="https://goproxy.cn,direct"
# 编译，在image内运行
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go

FROM zhugeqing/alpine
# 拷贝缓存区构建完成的文件（这样镜像还是非常小的）
COPY --from=builder /app/main ./main
# 暴露端口，只是显示而已
# EXPOSE 50001
# 运行 可执行文件
CMD [ "/main" ]