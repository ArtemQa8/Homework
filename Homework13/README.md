# Homework12: Шахматный REST API + gRPC

REST API + gRPC для игры в шахматы, консольные клиенты, Swagger-документация.

Реализовано:
- CRUD для игроков, игр, ходов (REST и gRPC)
- Выполнение ходов с проверкой правил шахмат
- Автоматические ходы и фоновые симуляции
- JWT-авторизация (логин/пароль из `.env`)
- Защита mutation-эндпоинтов (POST/PUT/DELETE) через Bearer-токен
- Swagger UI с описанием всех эндпоинтов
- Веб-страницы наблюдателя (`/` и `/game?id=N`)
- gRPC-сервер и клиент с тем же функционалом (Player/Game/Move)

## Требования

- Go 1.26+
- (опционально) `swag` CLI — если нужно перегенерировать Swagger
- (опционально) `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc` — если нужно перегенерировать gRPC-код

## Настройка

Все команды выполняются из папки `Homework12` (там, где лежит `go.mod`).

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

Из папки `Homework12`:

Сервер:

    go run ./cmd/server

Клиент (в отдельном терминале, тоже из `Homework12`):

    go run ./cmd/client

После запуска сервера доступны:

- Swagger UI:              http://localhost:8080/swagger/index.html
- Список игр:              http://localhost:8080/
- Наблюдение за игрой:     http://localhost:8080/game?id=<ID>

## Запуск gRPC

Помимо REST API, проект предоставляет gRPC-интерфейс с тем же функционалом.

### Генерация Go-кода из .proto (опционально)

Схема лежит в `internal/proto/chess.proto`. Если её меняли — перегенерировать:

    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/proto/chess.proto

Требуются установленные `protoc`, `protoc-gen-go`, `protoc-gen-go-grpc`.

### Запуск gRPC-сервера

Из папки `Homework12`:

    go run ./cmd/grpc-server

Сервер слушает порт 50051. Можно запускать одновременно с REST-сервером
(порт 8080)

> ⚠️ **Важно:** хранилище — общий JSON-файл без межпроцессной
> синхронизации. Не запускайте REST- и gRPC-серверы одновременно,
> если планируете записывать данные — изменения одного процесса
> могут быть потеряны при сохранении другим. Для тестирования
> по очереди: сначала один, потом другой.

### Запуск gRPC-клиента

В отдельном терминале:

    go run ./cmd/grpc-client

Клиент интерактивный, с меню:

- **Игроки** — CRUD через PlayerService
- **Игры** — CRUD + MakeMove + AutoMove через GameService
- **Ходы** — CRUD через MoveService

### Реализованные gRPC-сервисы

PlayerService (CRUD), GameService (CRUD + MakeMove + AutoMove), MoveService (CRUD).

## Проверка REST через Swagger

1. Откройте Swagger UI в браузере.
2. Выполните `POST /api/login` — в теле передайте логин и пароль из вашего `.env`:

       {
         "логин": "<ваш LOGIN>",
         "пароль": "<ваш PASSWORD>"
       }

3. В ответе придёт `{"токен": "eyJ..."}`. Скопируйте токен.
4. Нажмите кнопку **Authorize** (справа сверху).
5. В поле `Value` вставьте: `Bearer <ваш токен>` (слово Bearer, пробел, токен).
6. Нажмите Authorize → Close.

Теперь защищённые эндпоинты (POST/PUT/DELETE) работают — Swagger автоматически
добавляет заголовок `Authorization`.

GET-эндпоинты открыты без авторизации — для наблюдателей.

## Проверка REST через клиент

При старте клиент попросит логин и пароль. Введите те же данные, что в `.env`.
Клиент сам получит токен и будет прикладывать его ко всем изменяющим запросам.

Дальше — обычное меню: создать игру, подключиться к существующей, симуляции.

## Перегенерация Swagger (опционально)

Если меняли комментарии хендлеров в `cmd/server/handlers.go`:

    swag init -g cmd/server/main.go -o docs --parseDependency --parseInternal

## Структура проекта

### Код

**REST (HTTP):**

- `cmd/server`              — REST-сервер (main, handlers, middleware)
- `cmd/client`              — консольный REST-клиент

**gRPC:**

- `cmd/grpc-server`         — точка входа gRPC-сервера
- `cmd/grpc-client`         — консольный gRPC-клиент
- `internal/api/grpc`       — реализация gRPC-сервисов (player, game, move, convert)
- `internal/proto`          — `.proto` схема и сгенерированный Go-код

**Общее:**

- `internal/auth`           — генерация и проверка JWT
- `internal/config`         — загрузка `.env` в структуру Config
- `internal/model`          — доменные модели, правила шахмат
- `internal/repository`     — слой доступа к данным
- `internal/dto`            — DTO для REST API
- `internal/service`        — вспомогательная логика (случайный ход, логгер)

### Данные и конфигурация (в корне проекта)

- `.env`                    — ваши секреты (не коммитится)
- `.env.example`            — шаблон для `.env`
- `data/games.json`         — игры
- `data/moves.json`         — ходы
- `data/players.json`       — игроки
- `docs/`                   — сгенерированный Swagger (docs.go, swagger.json, swagger.yaml)