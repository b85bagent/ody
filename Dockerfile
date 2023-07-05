FROM golang:1.18.1

RUN cd /go/src
RUN mkdir -p /go/src/agent	
COPY ./ /go/src/agent	
RUN cd /go/src/agent && export GO111MODULE=on && go build -race
WORKDIR /go/src/agent	

CMD ["./agent"]