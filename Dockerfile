FROM golang:1.16.4-buster as builder
WORKDIR /go/src/github.com/yangjing0630/go-stream
#ENV GOPROXY=https://goproxy.cn,direct
COPY . .
RUN make build_linux

FROM debian:stretch-slim
COPY --from=builder /go/src/github.com/yangjing0630/go-stream/bin /lal/bin
COPY --from=builder /go/src/github.com/yangjing0630/go-stream/conf /lal/conf
