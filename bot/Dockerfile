FROM golang:1.19

WORKDIR /go/src

RUN apt-get update && apt-get install -y make

ENTRYPOINT [ "make" ]
