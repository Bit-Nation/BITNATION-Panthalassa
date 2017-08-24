FROM golang:1.8.3-stretch

#Copy project
COPY . /go/src/panthalassa

#Copy config
COPY ./.docker/panthalassa/config.json /go/src/panthalassa/panthalassa/config.json

WORKDIR /go/src/panthalassa

#Install dependencies
RUN go get -v -d ./panthalassa/

#Build executable
RUN go build -o /build/panthalassa panthalassa/panthalassa.go

ENTRYPOINT ["/build/panthalassa"]

EXPOSE 80