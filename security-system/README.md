# security-system

Три сервиса, собранные в одну Docker-сеть: прокси перед сервисом аутентификации,
плюс агент сбора системных логов со своим приёмником.

## Архитектура

```
Клиент ──▶ proxy:8888 ──▶ auth:8080 ──▶ postgres / redis
                                          (не доступны напрямую снаружи)

siem-agent ──▶ log-collector:8080 (тестовый приёмник логов, отдельно от auth)
```

**Почему auth не открыт наружу напрямую:** единственная точка входа в систему
для внешнего клиента — proxy. Если открыть порт auth наружу тоже, rate limiting
и security-проверки прокси можно тривиально обойти, обратившись к auth напрямую.

**Почему siem-agent шлёт логи не в auth:** это разные по смыслу сервисы — auth
занимается аутентификацией пользователей, а не приёмом системных логов ОС.
`log-collector` — простой тестовый Flask-сервер для демонстрации; в реальной
системе на его месте был бы Elasticsearch/Loki/аналог.

## Запуск

```bash
cp .env.example .env
# отредактируй .env — обязательно смени JWT_SECRET на случайную строку от 32 символов

make build
make up
make test
```

## Проверка вручную

```bash
# Регистрация через прокси
curl -X POST http://localhost:8888/register \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@mail.com","password":"test123"}'

# Логин через прокси — получаем access_token
curl -X POST http://localhost:8888/login \
  -H "Content-Type: application/json" \
  -d '{"username":"test","password":"test123"}'

# Validate через прокси (нужен Authorization — теперь пропускается корректно)
curl http://localhost:8888/validate \
  -H "Authorization: Bearer <access_token>"

# Что накопил приёмник логов от SIEM-агента
curl http://localhost:8081/api/logs
```

## Важно: какой docker-compose.yml использовать

Внутри `auth/` и `siem-agent/` остались их собственные `docker-compose.yml` —
это конфиги для запуска **каждого проекта отдельно и независимо** (как было
до объединения). Для запуска всей системы целиком используется **только**
`docker-compose.yml` в корне `security-system/` — команды `make build`/`make up`
запускаются из корня и используют именно его. Не пытайся одновременно поднять
корневой и вложенный compose — оба используют одинаковые порты (5432, 6379,
8080) и будут конфликтовать.

## Известные ограничения этой сборки

- **`SERVER_URL` для siem-agent задаётся только через `siem-agent/config.json`**,
  переменная окружения `SERVER_URL` в environment сервиса не читается — конфиг
  агента загружается только из JSON-файла (`config.Load`).
- **Rate limiting в proxy — по `r.RemoteAddr`, то есть по IP контейнера/клиента
  на уровне TCP-соединения**, а не по `X-Forwarded-For` (заголовок от клиента
  специально игнорируется — иначе rate limiting тривиально обходится подменой
  заголовка).
- Все ограничения отдельных сервисов (auth, siem-agent) из их собственных
  README остаются в силе — объединение в общую сеть их не устраняет.
