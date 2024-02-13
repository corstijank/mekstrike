FROM golang:1.21-alpine

WORKDIR /go/src/app
COPY . .
WORKDIR /go/src/app/gamemaster

RUN go get -d -v ./...
RUN go install -v ./...

CMD ["gamemaster"]