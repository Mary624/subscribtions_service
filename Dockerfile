FROM golang:1.26.1-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/subscriptions-service ./cmd/subscriptions/main.go

# Финальный образ
FROM alpine:3.18

RUN apk add --no-cache curl ca-certificates

WORKDIR /app

# Копируем бинарник и миграции
COPY --from=builder /app/subscriptions-service .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/config ./config

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD curl -f http://localhost:8080/health || exit 1

CMD ["./subscriptions-service"]