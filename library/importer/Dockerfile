FROM golang:1.17-alpine

RUN apk add --update \
    curl \
    && rm -rf /var/cache/apk/*

WORKDIR /go/src/app
COPY . .

WORKDIR /go/src/app/library/importer
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["/bin/sh", "-c", "importer && sleep 5 && curl -X POST localhost:3500/v1.0/shutdown"]