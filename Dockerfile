FROM golang:alpine
WORKDIR $GOPATH/src/application
ADD . ./
ENV GO111MODULE=on
ENV GOPROXY="https://goproxy.io"
RUN apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo 'Asia/Shanghai' >/etc/timezone \
    && apk del tzdata
RUN go build -o application .