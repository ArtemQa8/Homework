# Homework15: Шахматный REST API + gRPC (PostgreSQL)

REST API + gRPC для игры в шахматы, консольные клиенты, Swagger-документация, юнит-тесты бизнес-логики, реляционная БД и CI.

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
- **Хранение данных в PostgreSQL** (основная БД вместо MongoDB)
- **Миграции схемы через `goose`**
- **Типобезопасная генерация запросов через `sqlc`**
- **Транзакции** для атомарных операций (создание игры, выполнение хода)
- **Скрипты миграции** JSON -> PostgreSQL и экспорта PostgreSQL -> JSON

## Требования

- Go 1.26+
- **Docker Desktop** (или Docker + Docker Compose)
- (опционально) `swag` CLI — если нужно перегенерировать Swagger
- (опционально) `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc` — если нужно перегенерировать gRPC-код
- (опционально) `goose` и `sqlc` CLI — если нужно менять миграции или запросы

## Настройка

Все команды выполняются из папки `Homework15` (там, где лежит `go.mod`).

1. Скопируйте `.env.example` в `.env`:

       cp .env.example .env

2. Откройте `.env` и задайте свои значения:

       LOGIN=<ваш логин>
       PASSWORD=<ваш пароль>
       JWT_SECRET=<случайная строка минимум 32 символа>
       JWT_TTL=1h
       
       # Настройки PostgreSQL (основная БД)
       PG_DSN=postgres://chess:chess@localhost:5432/chess?sslmode=disable

Файл `.env` добавлен в `.gitignore` и не коммитится.

## Запуск баз данных (Docker)

Перед запуском серверов необходимо поднять PostgreSQL. В папке `Homework15` выполните:

    docker compose up -d

*(В `docker-compose.yml` также присутствуют сервисы MongoDB и Redis — они оставлены для обратной совместимости с HW14 и в HW15 не используются).*

## Миграции схемы

Схема БД описана в `db/migrations/`. Применить миграции:

    goose -dir db/migrations postgres "postgres://chess:chess@localhost:5432/chess?sslmode=disable" up

Или через `go run`:

    go run github.com/pressly/goose/v3/cmd/goose@latest -dir db/migrations postgres "postgres://chess:chess@localhost:5432/chess?sslmode=disable" up

## Генерация кода из SQL (sqlc)

Запросы лежат в `internal/repository/postgres/queries/`. После их изменения перегенерируйте Go-код:

    sqlc generate

## Миграция данных (JSON -> PostgreSQL)

Если у вас есть старые данные в JSON-файлах (папка `data`) от предыдущих версий проекта, вы можете перенести их в PostgreSQL. Скрипт сохранит все старые ID.

    go run ./cmd/import-from-json

## Экспорт данных (PostgreSQL -> JSON)

Для создания бэкапа в виде JSON-файлов:

    go run ./cmd/export-to-json

Данные сохраняются в папку `export_data/`.

## Запуск REST

Из папки `Homework15`:

Сервер:

    go run ./cmd/rest-server

Клиент (в отдельном терминале, тоже из `Homework15`):

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

Из папки `Homework15`:

    go run ./cmd/grpc-server

Сервер слушает порт 50051. Можно запускать одновременно с REST-сервером
(порт 8080). *Благодаря PostgreSQL параллельная работа серверов безопасна.*

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
- **Слой PostgreSQL** (`internal/repository/postgres`) — требует запущенного Docker-контейнера (интеграционные тесты).
- **Слой MongoDB** (`internal/repository/mongodb`) — оставлен от предыдущей версии проекта (HW14), в HW15 не используется.
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

- `cmd/import-from-json`     — скрипт миграции данных из JSON в PostgreSQL
- `cmd/export-to-json`       — скрипт экспорта данных из PostgreSQL в JSON
- `db/migrations`            — SQL-миграции для goose
- `docker-compose.yml`       — конфигурация PostgreSQL (MongoDB и Redis оставлены для совместимости)
- `internal/repository/filestore` — старое JSON-хранилище (для миграции)
- `internal/repository/mongodb`   — альтернативное хранилище (HW14, не используется)
- `internal/repository/postgres`  — основное хранилище на PostgreSQL
- `internal/repository/postgres/queries` — SQL-запросы для sqlc
- `internal/repository/postgres/sqlc`    — сгенерированный код

**Общее:**

- `internal/auth`            — генерация и проверка JWT
- `internal/config`          — загрузка `.env` в структуру Config
- `internal/model`           — доменные модели, правила шахмат
- `internal/service`         — бизнес-логика: Player, Game, Move
- `internal/dto`             — DTO для REST API
- `internal/render/tui`      — рендер доски для терминала
- `docs/`                    — сгенерированный Swagger