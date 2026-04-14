FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/api ./cmd/api
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/worker ./cmd/worker

FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /out/api /app/api
COPY --from=builder /out/worker /app/worker

EXPOSE 8080