FROM golang:1.8.3-stretch

COPY . /go/src/panthalassa

WORKDIR /go/src/panthalassa

RUN go get -v -d ./panthalassa/

EXPOSE 80