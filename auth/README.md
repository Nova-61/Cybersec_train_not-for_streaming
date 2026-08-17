# Auth Service (JWT + Redis)

Сервис аутентификации: регистрация, логин, refresh-токены, logout, хранение сессий в Redis.

## Функции

- Регистрация с bcrypt-хешированием пароля
- Логин (access + refresh JWT)
- Refresh с token rotation и проверкой сессии в Redis
- Logout (инвалидация сессии)
- Middleware для защищённых роутов
- Graceful shutdown

## Быстрый старт

### 1. Задать секреты

Перед запуском обязательно задай свой `JWT_SECRET` (минимум 32 символа) —
без него сервис не запустится. Не используй значение из docker-compose.yml
по умолчанию, оно только для примера.

Проще всего — создать `.env` файл рядом с `docker-compose.yml`:

```
JWT_SECRET=<сгенерируй случайную строку от 32 символов>
```

и подключить его в docker-compose.yml через `env_file: .env`
(сам файл `.env` должен быть в `.gitignore` и никогда не попадать в git).

### 2. Запуск через Docker Compose

```bash
docker-compose up -d
```

### 3. Запуск локально (без Docker для самого сервиса)

```bash
go mod download
docker-compose up -d postgres redis
JWT_SECRET=<твой секрет от 32 символов> go run ./cmd/auth
```

## API

### Регистрация

```bash
curl -X POST http://localhost:8080/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@mail.com","password":"secret123"}'
```

### Логин

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{"username":"john","password":"secret123"}'
```

### Refresh

```bash
curl -X POST http://localhost:8080/refresh \
  -H "Content-Type: application/json" \
  -d '{"refresh_token":"<refresh_token>"}'
```

### Валидация токена (защищённый роут)

```bash
curl -X GET http://localhost:8080/validate \
  -H "Authorization: Bearer <access_token>"
```

### Logout (защищённый роут)

```bash
curl -X POST http://localhost:8080/logout \
  -H "Authorization: Bearer <access_token>"
```

## Переменные окружения

| Переменная | По умолчанию | Обязательна | Описание |
|---|---|---|---|
| `JWT_SECRET` | — | да | Секрет для подписи JWT, минимум 32 символа |
| `PORT` | `8080` | нет | Порт сервера |
| `DB_HOST` | `localhost` | нет | PostgreSQL хост |
| `DB_PORT` | `5432` | нет | PostgreSQL порт |
| `DB_USER` | `postgres` | нет | PostgreSQL пользователь |
| `DB_PASSWORD` | `password` | нет | PostgreSQL пароль |
| `DB_NAME` | `auth_db` | нет | PostgreSQL база |
| `REDIS_ADDR` | `localhost:6379` | нет | Redis адрес |

## Известные ограничения (осознанные упрощения учебного проекта)

- **Logout не отзывает уже выданный access-токен**, только убивает возможность
  его обновить через `/refresh`. Access-токен остаётся валиден по подписи
  до истечения `exp` (по умолчанию 1 час). Полное решение — token revocation
  через denylist по `jti` или хранение refresh-токенов в БД со статусом.
- **Одна активная сессия на пользователя.** Логин с нового устройства
  перезаписывает сессию в Redis для предыдущего.
- **Таблицы создаются через `CREATE TABLE IF NOT EXISTS` при старте**,
  а не через полноценный migration-инструмент (`golang-migrate`, `goose`).
  Подходит для разработки, не для продакшена со сложной эволюцией схемы.
- **MFA не реализована** — несмотря на то, что упоминается в оригинальном
  плане проекта, в этом коде её нет.