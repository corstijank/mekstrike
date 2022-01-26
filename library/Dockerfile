FROM golang:1.16-alpine

WORKDIR /go/src/app
COPY . .

WORKDIR /go/src/app/library
RUN go get -d -v ./...
RUN go install -v ./...

CMD ["library"]