# 第一階段
FROM golang:1.18.1-alpine AS builder

RUN apk update && apk add --no-cache git

WORKDIR /app

#ENV GOPROXY=https://goproxy.io // go mod download 失敗請把註解拿掉

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o newProject

# 第二階段
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/newProject .
# COPY --from=builder /app/yaml ./yaml
# COPY --from=builder /app/blackbox_exporter ./blackbox_exporter

EXPOSE 8080

CMD ["./newProject"]