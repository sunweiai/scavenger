#FROM dev-scm-local.shijicloud.com/kunlun/golang:1.20.4-alpine AS builder
FROM nginx
WORKDIR /application
ADD . ./
ENV GO111MODULE=on
ENV GOPROXY="https://proxy.golang.com.cn,direct"
RUN go build -o scavenger main.go

#FROM alpine

#WORKDIR /application
#COPY --from=builder /application/scavenger /application/scavenger

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN echo 'Asia/Shanghai' >/etc/timezone
# 设置编码
ENV LANG C.UTF-8
CMD ["./scavenger"]