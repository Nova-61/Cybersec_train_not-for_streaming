# HTTP/HTTPS Proxy Server

## Описание
Прокси-сервер для фильтрации и логирования HTTP/HTTPS трафика.

## Функции
- ✅ Базовый ReverseProxy
- ✅ Логирование в JSON
- ✅ Фильтрация заголовков
- ✅ Rate Limiting (10 запросов в минуту)
- ✅ Graceful Shutdown
- ✅ Поддержка HTTPS (опционально)

## Запуск

### 1. Запустить тестовый сервер
```bash
go run test_server.go
```

### 2. Запустить прокси
```bash
go run ./cmd/proxy -listen :8888 -target http://localhost:8080
```

### 3. Проверить
```bash
curl http://localhost:8888/
curl http://localhost:8888/users
```

## Сборка
```bash
go build -o proxy ./cmd/proxy
```

## Тестирование
```bash
go test ./...
go test -cover ./...
```

## Безопасность
```bash
gosec ./...
```

## Структура проекта
```
proxy/
├── cmd/proxy/main.go      # Точка входа
├── internal/
│   ├── config/            # Конфигурация
│   ├── handler/           # Логика прокси
│   └── middleware/        # Middleware (logger, security, rate_limit)
├── go.mod
├── test_server.go
└── README.md
```
