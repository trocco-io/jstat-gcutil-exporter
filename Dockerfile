FROM golang:1.13.4

WORKDIR /go/src/app

RUN go get -u github.com/golang/dep/cmd/dep

COPY . .
