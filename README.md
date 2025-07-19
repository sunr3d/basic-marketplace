# Basic Marketplace

Тестовое задание для VK Tech, Backend Go developer.

## Как запустить

```bash
docker-compose up --build
```

## Технологии

- Go, Gin, GORM
- PostgreSQL, Redis
- JWT, Docker, Docker Compose

## Структура

- cmd/ — точка входа
- internal/ — бизнес-логика, обработчики, инфраструктура
- models/ — основные структуры

## Описание API

TODO

## Переменные окружения

Для настройки приложения можно использовать переменные окружения или файл .env (см. .env.example):

- `HTTP_HOST` — адрес, на котором слушает сервер (по умолчанию 0.0.0.0)
- `HTTP_PORT` — порт сервера (по умолчанию 8080)
- `LOG_LEVEL` — уровень логирования (debug, info, warn, error; по умолчанию debug)

**PostgreSQL:**

- `POSTGRES_HOST` — адрес БД (по умолчанию localhost)
- `POSTGRES_PORT` — порт БД (по умолчанию 5432)
- `POSTGRES_USER` — пользователь БД (по умолчанию postgres)
- `POSTGRES_PASSWORD` — пароль БД (по умолчанию postgres)
- `POSTGRES_DATABASE` — имя БД (по умолчанию marketplace)

**Redis:**

- `REDIS_ADDR` — адрес Redis (по умолчанию localhost:6379)
- `REDIS_PASSWORD` — пароль Redis (по умолчанию пусто)
- `REDIS_DB` — номер базы Redis (по умолчанию 0)

Если файл .env отсутствует, используются значения по умолчанию.
