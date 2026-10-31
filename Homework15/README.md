# Homework14: Шахматный REST API + gRPC (MongoDB + Redis)

REST API + gRPC для игры в шахматы, консольные клиенты, Swagger-документация, юнит-тесты бизнес-логики, базы данных NoSQL и CI.

## Что реализовано

- CRUD для игроков, игр, ходов (REST и gRPC)
- Выполнение ходов с проверкой правил шахмат
- Автоматические ходы и фоновые симуляции
- JWT-авторизация (логин/пароль из `.env`)
- Защита mutation-эндпоинтов (POST/PUT/DELETE) через Bearer-токен
- Swagger UI с описанием всех эндпоинтов
- Веб-страницы наблюдателя 
- Юнит-тесты бизнес-логики (black box)
- CI: сборка и прогон тестов на каждый push/PR
- **Хранение данных в MongoDB** (вместо JSON-файлов)
- **Логирование изменений в Redis** с временем жизни (TTL)
- **Docker Compose** для запуска MongoDB и Redis
- **Скрипт миграции** старых данных из JSON в MongoDB

## Требования

- Go 1.26+
- **Docker Desktop** (или Docker + Docker Compose)
- (опционально) `swag` CLI — если нужно перегенерировать Swagger
- (опционально) `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc` — если нужно перегенерировать gRPC-код

## Настройка

Все команды выполняются из папки `Homework14` (там, где лежит `go.mod`).

1. Скопируйте `.env.example` в `.env`:

       cp .env.example .env

2. Откройте `.env` и задайте свои значения:

       LOGIN=<ваш логин>
       PASSWORD=<ваш пароль>
       JWT_SECRET=<случайная строка минимум 32 символа>
       JWT_TTL=1h
       
       # Настройки баз данных
       MONGO_URI=mongodb://localhost:27017
       MONGO_DB=chess
       REDIS_ADDR=localhost:6379
       REDIS_TTL=24h

Файл `.env` добавлен в `.gitignore` и не коммитится.

## Запуск баз данных (Docker)

Перед запуском серверов необходимо поднять MongoDB и Redis. В папке `Homework14` выполните:

    docker compose up -d

*(Флаг `-d` запускает контейнеры в фоновом режиме).*

## Миграция данных (опционально)

Если у вас есть старые данные в JSON-файлах (папка `data`) от предыдущих версий проекта, вы можете перенести их в MongoDB. Скрипт сохранит все старые ID.

    go run ./cmd/migrate

## Запуск REST

Из папки `Homework14`:

Сервер:

    go run ./cmd/rest-server

Клиент (в отдельном терминале, тоже из `Homework14`):

    go run ./cmd/rest-client

После запуска сервера доступны:

- Swagger UI:              http://localhost:8080/swagger/index.html
- Список игр:              http://localhost:8080/
- Наблюдение за игрой:     http://localhost:8080/game?id=<ID>

## Запуск gRPC

Помимо REST API, проект предоставляет gRPC-интерфейс с тем же функционалом.

### Генерация Go-кода из .proto (опционально)

Схема лежит в `internal/proto/chess.proto`. Если её меняли — перегенерировать:

    protoc --go_out=. --go_opt=paths=source_relative \
        --go-grpc_out=. --go-grpc_opt=paths=source_relative \
        internal/proto/chess.proto

### Запуск gRPC-сервера

Из папки `Homework14`:

    go run ./cmd/grpc-server

Сервер слушает порт 50051. Можно запускать одновременно с REST-сервером
(порт 8080). *Благодаря MongoDB параллельная работа серверов безопасна.*

### Запуск gRPC-клиента

В отдельном терминале:

    go run ./cmd/grpc-client

Клиент интерактивный, с меню:
- **Игроки** — CRUD через PlayerService
- **Игры** — CRUD + MakeMove + AutoMove через GameService
- **Ходы** — CRUD через MoveService

## Тесты

Все тесты — **black box** (используют только публичный API пакетов).
Запускаются через `go test`.

### Запуск всех тестов

    go test ./...

### Запуск тестов конкретного пакета

    go test ./internal/model/
    go test ./internal/service/
    go test ./internal/repository/filestore/
    go test ./internal/config/

### Без кеша (реальный прогон)

    go test ./... -count=1

## Покрытие

### Текущее покрытие бизнес-логики

| Пакет        | Покрытие |
| ------------ | -------- |
| `model`      | ~80%     |
| `service`    | ~92%     |
| `filestore`  | ~92%     |
| `config`     | 100%     |
| `auth`       | ~88%     |

**Что не покрыто осознанно:**
- **Слой MongoDB** (`internal/repository/mongodb`) — требует запущенного Docker-контейнера (интеграционные тесты).
- **Рендер** (`internal/render/tui`, HTML-страницы) — визуальная часть.
- **Транспортный слой** (`internal/api/rest`, `internal/api/grpc`) — тестируется косвенно.

## CI

Проект использует **GitHub Actions**. Workflow лежит в
`.github/workflows/ci.yml` в корне репозитория.

**Что делает CI на каждый push и pull request:**

1. Скачивает код
2. Устанавливает Go 1.26
3. Скачивает зависимости (`go mod download`)
4. Собирает проект (`go build ./...`)
5. Запускает тесты (`go test ./... -count=1`)

## Структура проекта

### Код

**REST (HTTP):**

- `cmd/rest-server`          — REST-сервер (main, роуты, graceful shutdown)
- `cmd/rest-client`          — консольный REST-клиент
- `internal/api/rest`        — хендлеры и middleware

**gRPC:**

- `cmd/grpc-server`          — точка входа gRPC-сервера
- `cmd/grpc-client`          — консольный gRPC-клиент
- `internal/api/grpc`        — реализация gRPC-сервисов
- `internal/proto`           — `.proto` схема и сгенерированный Go-код

**Базы данных и утилиты:**

- `cmd/migrate`              — скрипт миграции данных из JSON в MongoDB
- `docker-compose.yml`       — конфигурация MongoDB и Redis
- `internal/repository/filestore` — старое JSON-хранилище (для миграции)
- `internal/repository/mongodb`   — новое хранилище на MongoDB и Redis

**Общее:**

- `internal/auth`            — генерация и проверка JWT
- `internal/config`          — загрузка `.env` в структуру Config
- `internal/model`           — доменные модели, правила шахмат
- `internal/service`         — бизнес-логика: Player, Game, Move
- `internal/dto`             — DTO для REST API
- `internal/render/tui`      — рендер доски для терминала
- `docs/`                    — сгенерированный Swagger