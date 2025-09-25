FROM golang:tip-alpine3.22

RUN go install github.com/air-verse/air@v1.63.0 && \
    go install github.com/swaggo/swag/cmd/swag@v1.16.6

WORKDIR /app