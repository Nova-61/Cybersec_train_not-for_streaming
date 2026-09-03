.PHONY: help build up down logs test clean

help:
	@echo "make build  - собрать все образы"
	@echo "make up     - поднять всю систему"
	@echo "make down   - остановить всё"
	@echo "make logs   - логи всех сервисов"
	@echo "make test   - прогнать базовую проверку через прокси"
	@echo "make clean  - остановить и удалить volumes (данные БД/Redis)"

build:
	docker compose build

up:
	docker compose up -d
	@echo "Proxy:        http://localhost:8888"
	@echo "Log-collector: http://localhost:8081/api/logs"

down:
	docker compose down

logs:
	docker compose logs -f

test:
	@echo "=== Регистрация через прокси ==="
	curl -s -X POST http://localhost:8888/register \
	  -H "Content-Type: application/json" \
	  -d '{"username":"test","email":"test@mail.com","password":"test123"}'; echo
	@echo "=== Логин через прокси ==="
	curl -s -X POST http://localhost:8888/login \
	  -H "Content-Type: application/json" \
	  -d '{"username":"test","password":"test123"}'; echo

clean:
	docker compose down -v
