FROM golang:1.18-alpine

RUN apk add git
WORKDIR /usr/src/app

COPY . .

RUN go mod download && go mod verify

RUN go build -v -o /usr/local/bin/app ./cmd/serve/main.go

ENTRYPOINT ["app"]