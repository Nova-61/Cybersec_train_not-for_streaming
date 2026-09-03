# Сборка
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /auth ./cmd/auth

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget && \
    adduser -D -H appuser

WORKDIR /home/appuser

COPY --from=builder /auth .
RUN chown appuser:appuser /auth

USER appuser

EXPOSE 8080

ENTRYPOINT ["./auth"]
