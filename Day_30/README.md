# День 30 — Нагрузочное тестирование API

## Что здесь есть

- `main.go` — простой REST API с тремя эндпоинтами:
  - `GET /users`
  - `GET /users/1`
  - `POST /users`
- `vegeta-targets.txt` — примеры запросов для `vegeta`
- `wrk.sh` — пример команды для `wrk`

## Запуск API

Откройте терминал в `Day_30` и выполните:

```bash
go run main.go
```

Сервер будет доступен на `http://localhost:8080`.

## Тестирование vegeta

1. Создайте файл `vegeta-targets.txt`:

```text
GET http://localhost:8080/users
GET http://localhost:8080/users/1
POST http://localhost:8080/users
Content-Type: application/json

{"name":"Anna","email":"anna@mail.com","age":28}
```

2. Запустите тест:

```bash
echo "GET http://localhost:8080/users" | vegeta attack -duration=10s -rate=100 -output=results.bin
vegeta report -type=text results.bin
```

3. Отчёт:

```bash
vegeta report results.bin
vegeta report -type=text results.bin
vegeta plot results.bin > plot.html
```

## Тестирование wrk

Пример команды:

```bash
wrk -t12 -c400 -d30s http://localhost:8080/users
```

## Задача 3: 500 RPS

```bash
echo "GET http://localhost:8080/users" | vegeta attack -duration=60s -rate=500 -output=results.bin
vegeta report -type=text results.bin
```

Если сервер не выдержит, можно:

- уменьшить логирование
- убрать лишние операции в обработчиках
- использовать пул соединений для реальной БД
- добавить индексы в таблицы
