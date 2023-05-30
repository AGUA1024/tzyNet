FROM golang:1.20.2-alpine

# 创建文件夹
RUN mkdir /app

# 设置工作目录
WORKDIR /app

#将Dockerfile 中的文件存储到/app下
ADD . /app

# 因为已经是在 /app下了，所以使用  ./
RUN go mod download
RUN go mod tidy
RUN go build -o main ./gateway.go

# 暴露的端口
EXPOSE 8000
EXPOSE 8001

#设置容器的启动命令，CMD是设置容器的启动指令
CMD /app/main -tag ${NODEMARK}
