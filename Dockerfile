FROM golang:1.18.1

WORKDIR /go/src/remote_write

COPY . .

RUN export GO111MODULE=on && go build -race

EXPOSE 8080

CMD ["./remote_write"]
