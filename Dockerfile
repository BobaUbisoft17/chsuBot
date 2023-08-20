FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY internal ./internal
COPY pkg ./pkg
COPY cmd/app ./cmd/app

RUN go build -o /app/chsuBot ./cmd/app/

CMD ["/app/chsuBot"]