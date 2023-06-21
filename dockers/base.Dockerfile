FROM golang:1.20.4-buster AS builder

# 配置 Go SDK 镜像源
ENV GOPROXY=https://goproxy.cn,direct

# 安装常用的 Go 工具
RUN	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.30.0
RUN	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
RUN	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-errors/v2@v2.0.0-20230515030202-6d741828c2d4
RUN	go install github.com/go-kratos/kratos/cmd/kratos/v2@v2.0.0-20230515030202-6d741828c2d4
RUN	go install github.com/go-kratos/kratos/cmd/protoc-gen-go-http/v2@v2.0.0-20230515030202-6d741828c2d4
RUN	go install github.com/google/gnostic/cmd/protoc-gen-openapi@v0.6.8
RUN	go install github.com/envoyproxy/protoc-gen-validate@v1.0.1
RUN	go install github.com/google/wire/cmd/wire@v0.5.0
RUN go install github.com/bufbuild/buf/cmd/buf@v1.21.0

# # 安装 curl
# RUN apk add curl
# RUN apk add make

RUN apt-get update && apt-get install -y unzip 
# 安装 protoc 和 protoc-gen-go
RUN curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v22.5/protoc-22.5-linux-x86_64.zip
RUN unzip protoc-22.5-linux-x86_64.zip -d /usr/local
RUN rm protoc-22.5-linux-x86_64.zip