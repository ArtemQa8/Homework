# Homework13: Шахматный REST API + gRPC

REST API + gRPC для игры в шахматы, консольные клиенты, Swagger-документация,
юнит-тесты бизнес-логики и CI.

## Что реализовано

- CRUD для игроков, игр, ходов (REST и gRPC)
- Выполнение ходов с проверкой правил шахмат
- Автоматические ходы и фоновые симуляции
- JWT-авторизация (логин/пароль из `.env`)
- Защита mutation-эндпоинтов (POST/PUT/DELETE) через Bearer-токен
- Swagger UI с описанием всех эндпоинтов
- Веб-страницы наблюдателя (`/` и `/game?id=N`)
- gRPC-сервер и клиент с тем же функционалом
- **Юнит-тесты бизнес-логики** (black box)
- **CI**: сборка и прогон тестов на каждый push/PR

## Требования

- Go 1.26+
- (опционально) `swag` CLI — если нужно перегенерировать Swagger
- (опционально) `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc` — если нужно
  перегенерировать gRPC-код

## Настройка

Все команды выполняются из папки `Homework13` (там, где лежит `go.mod`).

1. Скопируйте `.env.example` в `.env`:

       cp .env.example .env

2. Откройте `.env` и задайте свои значения:

       LOGIN=<ваш логин>
       PASSWORD=<ваш пароль>
       JWT_SECRET=<случайная строка минимум 32 символа>
       JWT_TTL=1h

   Эти данные нужны для входа в API. Регистрация не предусмотрена —
   логин/пароль хранятся только в `.env` на сервере.

Файл `.env` добавлен в `.gitignore` и не коммитится.

## Запуск REST

Из папки `Homework13`:

Сервер:

    go run ./cmd/rest-server

Клиент (в отдельном терминале, тоже из `Homework13`):

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

Требуются установленные `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`.

### Запуск gRPC-сервера

Из папки `Homework13`:

    go run ./cmd/grpc-server

Сервер слушает порт 50051. Можно запускать одновременно с REST-сервером
(порт 8080).

> ⚠️ **Важно:** хранилище — общий JSON-файл без межпроцессной
> синхронизации. Не запускайте REST- и gRPC-серверы одновременно,
> если планируете записывать данные — изменения одного процесса
> могут быть потеряны при сохранении другим.

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

### Запуск с подробным выводом

    go test -v ./...

### Запуск тестов конкретного пакета

    go test ./internal/model/
    go test ./internal/service/
    go test ./internal/repository/
    go test ./internal/config/

### Без кеша (реальный прогон)

    go test ./... -count=1

## Покрытие

### Проценты по каждому пакету

    go test -cover ./...

### Подробный отчёт по функциям

    go test -coverprofile coverage.out ./...
    go tool cover -func coverage.out

### HTML-отчёт (в браузере)

    go test -coverprofile coverage.out ./...
    go tool cover -html coverage.out

### Текущее покрытие бизнес-логики

| Пакет        | Покрытие |
| ------------ | -------- |
| `model`      | ~80%     |
| `service`    | ~92%     |
| `repository` | ~92%     |
| `config`     | 100%     |
| `auth`       | ~88%     |

**Что покрыто осознанно:** бизнес-логика, JSON-сериализация, HTTP-API
через сервисы, работа с хранилищем, валидация конфигурации.

**Что не покрыто осознанно:**
- **Рендер** (`internal/render/tui`, HTML-страницы) — визуальная часть,
  проверяется вручную
- **Транспортный слой** (`internal/api/rest`, `internal/api/grpc`) —
  тестируется косвенно через сервисы
- **DTO и сгенерированный код** (`internal/dto`, `internal/proto`, `docs/`) —
  только структуры или автогенерация
- **Точки входа** (`cmd/*`) — `main`-функции, запускаются вручную

## CI

Проект использует **GitHub Actions**. Workflow лежит в
`.github/workflows/ci.yml` в корне репозитория.

**Что делает CI на каждый push и pull request:**

1. Скачивает код
2. Устанавливает Go 1.26
3. Скачивает зависимости (`go mod download`)
4. Собирает проект (`go build ./...`)
5. Запускает тесты (`go test ./... -count=1`)

Статус CI виден во вкладке **Actions** репозитория и в статусе PR.

## Проверка REST через Swagger

1. Откройте Swagger UI в браузере.
2. Выполните `POST /api/login` — в теле передайте логин и пароль из `.env`:

       {
         "логин": "<ваш LOGIN>",
         "пароль": "<ваш PASSWORD>"
       }

3. В ответе придёт `{"токен": "eyJ..."}`. Скопируйте токен.
4. Нажмите кнопку **Authorize** (справа сверху).
5. В поле `Value` вставьте: `Bearer <ваш токен>`.
6. Нажмите Authorize → Close.

Теперь защищённые эндпоинты работают — Swagger автоматически
добавляет заголовок `Authorization`.

GET-эндпоинты открыты без авторизации — для наблюдателей.

## Проверка REST через клиент

При старте клиент попросит логин и пароль. Введите те же данные, что в `.env`.
Клиент сам получит токен и будет прикладывать его ко всем изменяющим
запросам.

Дальше — обычное меню: создать игру, подключиться к существующей,
симуляции.

### Команды в игровой сессии

- `автоход N` — сделать N автоматических ходов
- `симуляция N` — запустить N фоновых партий
- `стоп` — остановить симуляции
- `help` — список команд
- `выход` / `exit` — выйти из игры

## Перегенерация Swagger (опционально)

Если меняли комментарии хендлеров в `internal/api/rest/handlers.go`:

    swag init -g cmd/rest-server/main.go -o docs \
        --parseDependency --parseInternal

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

**Общее:**

- `internal/auth`            — генерация и проверка JWT
- `internal/config`          — загрузка `.env` в структуру Config
- `internal/model`           — доменные модели, правила шахмат
- `internal/repository`      — слой доступа к данным (JSON-хранилище)
- `internal/service`         — бизнес-логика: Player, Game, Move
- `internal/dto`             — DTO для REST API
- `internal/render/tui`      — рендер доски для терминала
- `docs/`                    — сгенерированный Swagger


### Данные и конфигурация (в корне `Homework13/`)

- `.env`                     — ваши секреты (не коммитится)
- `.env.example`             — шаблон для `.env`
- `data/games.json`          — игры
- `data/moves.json`          — ходы
- `data/players.json`        — игроки
- `data/counters.json`       — счётчики ID
- `.github/workflows/ci.yml` — CI-пайплайн (в корне репозитория)

### Тесты

Тесты лежат рядом с кодом (`*_test.go`). Помощники:

- `internal/model/testhelpers_test.go` — `emptyBoard`
- `internal/service/testhelpers_test.go` — `createTestGame`, `emptyServiceBoard`
- `internal/service/mock_test.go` — in-memory мок `Storage`