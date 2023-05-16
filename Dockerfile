FROM registry.cn-hangzhou.aliyuncs.com/qjwwy/go-builder:v1.20.4-bc23939 AS builder

COPY . /src

WORKDIR /src

ENV GOPROXY=https://goproxy.cn

RUN make release

FROM alpine:3.13.5

COPY --from=builder /src/bin /app

WORKDIR /app

ENV TZ=Asia/Shanghai
ENV GOPROXY=https://goproxy.cn

CMD ["./server", "-conf", "/app/configs/"]