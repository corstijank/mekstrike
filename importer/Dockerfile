FROM golang:1.21-alpine

WORKDIR /go/src/app
COPY . .
WORKDIR /go/src/app/importer

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["/bin/sh", "-c", "importer && sleep 5 && curl -X POST localhost:3500/v1.0/shutdown"]