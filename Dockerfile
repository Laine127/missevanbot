FROM golang:1.17-alpine
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories
RUN apk add --no-cache go
WORKDIR /app/missevan/
COPY . .
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct
RUN go build -ldflags "-s -w" -o missevan-bot cmd/main.go
CMD ./missevan-bot