# Basic Marketplace

## Описание

Backend REST API для маркетплейса объявлений.  
Реализовано на Go с использованием Clean Architecture, Gin, GORM, JWT, bcrypt, zap, Docker.

---

## Быстрый старт

### 1. Запуск через Docker

```bash
docker-compose up --build
```

### 2. Локальный запуск

```bash
go run cmd/main.go
```

### 3. Запуск через Make

```bash
make mocks # Сгенерировать моки под интерфейсы
```
```bash
make test # Запустить все тесты
```
```bash
make up # Поднять приложение и БД через docker-compose
```
```bash
make down # Остановить и удалить контейнеры, и зависимости
```
```bash
make clean # Удалить сгенерированные моки
```

---

## Переменные окружения

| Переменная        | Описание              | Пример значения |
| ----------------- | --------------------- | --------------- |
| POSTGRES_HOST     | Адрес Postgres        | localhost       |
| POSTGRES_PORT     | Порт Postgres         | 5432            |
| POSTGRES_USER     | Пользователь Postgres | postgres        |
| POSTGRES_PASSWORD | Пароль Postgres       | postgres        |
| POSTGRES_DB       | Имя БД                | marketplace     |
| JWT_SECRET        | Секрет для JWT        | supersecret     |
| HTTP_HOST         | Адрес HTTP-сервера    | 0.0.0.0         |
| HTTP_PORT         | Порт HTTP-сервера     | 8080            |

---

## Архитектура

- **Clean Architecture**:
  - `models/` — структуры БД
  - `interfaces/` — интерфейсы сервисов и репозиториев
  - `infra/` — инфраструктурные реализации (Postgres, репозитории)
  - `logic/` — бизнес-логика
  - `handlers/` — HTTP-ручки
  - `bootstrap/` — DI-контейнер
  - `entrypoint/` — сборка и запуск приложения
- **DI** через контейнер
- **GORM** для работы с Postgres
- **Gin** для HTTP API
- **zap** для логирования
- **testify, mockery** для тестов

---

## API

### POST /register

Регистрация пользователя

**Request:**

```json
{
  "login": "user1",
  "password": "Password123!"
}
```

**Response:**

```json
{
  "id": 1,
  "login": "user1",
  "created_at": "2025-07-21T12:00:00Z"
}
```

---

### POST /login

Авторизация, получение JWT

**Request:**

```json
{
  "login": "user1",
  "password": "Password123!"
}
```

**Response:**

```json
{
  "jwt": "jwt_token"
}
```

---

### POST /ads/create

Создание объявления (требует авторизации)

**Request:**

```json
{
  "title": "Продам велосипед",
  "description": "Почти новый",
  "image_url": "https://example.com/image.jpg",
  "price": 10000
}
```

**Response:**

```json
{
  "id": 1,
  "title": "Продам велосипед",
  "description": "Почти новый",
  "image_url": "https://example.com/image.jpg",
  "price": 10000,
  "owner_id": 1,
  "created_at": "2025-07-21T12:00:00Z"
}
```

---

### GET /ads

Лента объявлений (публичная)

**Query params:**

- `limit` — количество объявлений на страницу (по умолчанию 20, максимум 100)
- `offset` — смещение (начиная с какого объявления, например, offset=20 — с 21-го)
- `page` — номер страницы (альтернатива offset: если передан page — offset игнорируется; page=1 — первая страница)
- `min_price` — минимальная цена
- `max_price` — максимальная цена
- `sort_by` — поле сортировки (`created_at` или `price`)
- `order` — порядок сортировки (`asc` или `desc`)

**Примеры:**

- `/ads?limit=10&page=2` — вторая страница по 10 объявлений (offset=10)
- `/ads?limit=5&offset=15` — начиная с 16-го объявления, 5 штук
- `/ads?min_price=1000&max_price=5000&sort_by=price&order=asc&page=3` — фильтрация и пагинация вместе

**Response:**

```json
[
  {
    "id": 1,
    "title": "Продам велосипед",
    "description": "Почти новый",
    "image_url": "https://example.com/image.jpg",
    "price": 10000,
    "owner_login": "user1",
    "is_owner": false,
    "created_at": "2025-07-21 12:00:00"
  }
]
```

---

## Тесты

```bash
go test ./...
```

---

## Технологии

- Go 1.21+
- Gin
- GORM
- PostgreSQL
- JWT (github.com/golang-jwt/jwt/v5)
- bcrypt (golang.org/x/crypto)
- zap (go.uber.org/zap)
- testify, mockery
- Docker, docker-compose

---

## Структура проекта

```
cmd/                # main.go — точка входа
internal/
  bootstrap/        # DI-контейнер
  config/           # конфиг и загрузка из env
  infra/            # инфраструктурные реализации (Postgres, репозитории)
  interfaces/       # интерфейсы сервисов и репозиториев
  logic/            # бизнес-логика
  handlers/         # HTTP-ручки
  server/           # запуск HTTP-сервера с graceful shutdown
  entrypoint/       # сборка и запуск приложения
models/             # структуры БД
```

---

## Автор

[github.com/sunr3d](https://github.com/sunr3d)

---
