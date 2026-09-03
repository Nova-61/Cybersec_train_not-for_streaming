# Сборка
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /proxy ./cmd/proxy

# Финальный образ
FROM alpine:latest

RUN apk --no-cache add ca-certificates wget && \
    adduser -D -H appuser

WORKDIR /home/appuser

COPY --from=builder /proxy .
RUN chown appuser:appuser /proxy

USER appuser

EXPOSE 8888

ENTRYPOINT ["./proxy"]
