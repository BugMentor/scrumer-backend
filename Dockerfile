# syntax=docker/dockerfile:1

# Main application builder stage
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY *.go ./
COPY graph ./graph
COPY models ./models

RUN go build -o /scrumer-backend

# Seed application builder stage
FROM golang:1.25-alpine AS seed-builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/seed ./cmd/seed
COPY graph ./graph
COPY models ./models

RUN go build -o /scrumer-backend-seed ./cmd/seed

# Purge application builder stage
FROM golang:1.25-alpine AS purge-builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY cmd/purge ./cmd/purge
COPY graph ./graph
COPY models ./models

RUN go build -o /scrumer-backend-purge ./cmd/purge

# Final application image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /scrumer-backend ./scrumer-backend

EXPOSE 8080

CMD [ "/app/scrumer-backend" ]
