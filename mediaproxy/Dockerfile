FROM golang:1.17-alpine

WORKDIR /go/src/app
COPY . .

WORKDIR /go/src/app/mediaproxy
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["mediaproxy"]