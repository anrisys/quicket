FROM golang:1.24

RUN go install github.com/cosmtrek/air@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest

WORKDIR /app