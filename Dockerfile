FROM golang:1.18.1

WORKDIR /go/src/newProject

COPY . .

RUN export GO111MODULE=on && go build -race

EXPOSE 8080

CMD ["./newProject"]
