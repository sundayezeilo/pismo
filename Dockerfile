FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download
RUN go mod verify

COPY ./ ./

RUN go build -ldflags="-w -s" -o bin ./cmd

FROM alpine

COPY --from=builder /app/bin /usr/local/bin/server
COPY .env .env

ENTRYPOINT ["/usr/local/bin/server"]
