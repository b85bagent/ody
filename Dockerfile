FROM golang:1.18.1

RUN cd /go/src
RUN mkdir -p /go/src/Agent	
COPY ./ /go/src/Agent	
RUN cd /go/src/Agent && export GO111MODULE=on && go build -race
WORKDIR /go/src/Agent	

CMD ["./Agent"]