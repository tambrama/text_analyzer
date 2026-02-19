
## Описание проекта

Проект реализует два простых микросервиса на Go:
- **Сервис A (Text Receiver)**: принимает HTTP‑запросы с текстом, генерирует уникальный ID, отправляет текст в сервис B для анализа и хранит статус обработки в памяти.
- **Сервис B (Text Analyzer)**: принимает текст от сервиса A, считает статистику (количество слов, символов, предложений, среднюю длину слова) и синхронно возвращает результат.

Сервисы общаются по HTTP, поддерживают health‑checks, логирование, таймауты при сетевом запросе и корректное завершение (graceful shutdown). Хранение данных реализовано in‑memory в сервисе A.

## Архитектура

- **Язык**: Go 1.22
- **Сервис A (Receiver)**:
  - REST API:
    - `POST /api/v1/text` — принимает текст для анализа, возвращает `{id, status}`.
    - `GET /api/v1/status/{id}` — возвращает полное состояние запроса (статус, результат анализа).
    - `GET /api/v1/health` — health‑check.
  - При получении текста:
    - валидирует входные данные;
    - генерирует уникальный ID;
    - сохраняет запрос в in‑memory хранилище со статусом `PENDING`;
    - асинхронно (в отдельной горутине) отправляет текст в сервис B по HTTP с таймаутом;
    - по ответу сервиса B обновляет запись в хранилище: статус `DONE` или `FAILED`, сохраняет статистику / ошибку.

- **Сервис B (Analyzer)**:
  - REST API:
    - `POST /api/v1/analyze` — принимает `{id, text}`, возвращает статистику текста.
    - `GET /api/v1/health` — health‑check.
  - Считает:
    - количество слов (разделение по пробелам и знакам пунктуации);
    - количество символов (по рунам);
    - количество предложений (разделители `.`, `!`, `?`);
    - среднюю длину слова (по буквенным символам).

- **Взаимодействие**:
  - Сервис A вызывает `POST analyzer/api/v1/analyze` с JSON `{id, text}`.
  - Сервис B возвращает JSON со статистикой.
  - Сервис A сохраняет результат и делает его доступным по `GET /api/v1/status/{id}`.

## Структура проекта

- `cmd/app` — бинарь сервиса A (Text Receiver).
- `cmd/analyzer` — бинарь сервиса B (Text Analyzer).
- `internal/model` — модели домена (`TextRequest`, `TextStats`, статусы запроса).
- `internal/repository` — in‑memory репозиторий запросов (сервис A).
- `internal/service` — бизнес‑логика:
  - `ReceiverService` — валидация, генерация ID, вызов сервиса B, обновление статуса.
  - `AnalyzerService` — вычисление статистики текста.
- `internal/web` — HTTP‑обработчики и роутеры для обоих сервисов.
- `Dockerfile.receiver`, `Dockerfile.analyzer` — Docker‑образы для каждого сервиса.
- `docker-compose.yml` — запуск обоих сервисов вместе.

## Инструкция по запуску

### Локально (без Docker)

В требованиях предполагается Go 1.22.

1. **Собрать и запустить сервис B (Analyzer)**:
   ```bash
   go run ./cmd/analyzer
   ```
   По умолчанию он слушает `:8081` и имеет эндпоинты:
   - `POST http://localhost:8081/api/v1/analyze`
   - `GET http://localhost:8081/api/v1/health`

2. **Собрать и запустить сервис A (Receiver)**:
   В другом терминале:
   ```bash
   set ANALYZER_URL=http://localhost:8081  # Windows PowerShell/cmd
   go run ./cmd/app
   ```
   По умолчанию он слушает `:8080` и имеет эндпоинты:
   - `POST http://localhost:8080/api/v1/text`
   - `GET http://localhost:8080/api/v1/status/{id}`
   - `GET http://localhost:8080/api/v1/health`

Переменные окружения:
- `ANALYZER_URL` — URL сервиса B для сервиса A (по умолчанию `http://localhost:8081`).
- `RECEIVER_ADDR` — адрес прослушивания сервиса A (по умолчанию `:8080`).
- `ANALYZER_ADDR` — адрес прослушивания сервиса B (по умолчанию `:8081`).

### Запуск через Docker Compose

Необходимо установить Docker и docker‑compose.

1. Собрать и запустить:
   ```bash
   docker-compose up --build
   ```
2. После запуска:
   - Сервис A (receiver): `http://localhost:8080`
   - Сервис B (analyzer): `http://localhost:8081`

В `docker-compose.yml` уже настроены:
- отдельные контейнеры для каждого сервиса;
- проброс портов `8080` и `8081` на хост;
- переменная `ANALYZER_URL` для сервиса A (`http://analyzer:8081`).

## Примеры использования

### 1. Отправка текста на анализ

Запрос:
```bash
curl -X POST "http://localhost:8080/api/v1/text" ^
  -H "Content-Type: application/json" ^
  -d "{\"text\":\"Hello world. This is a test!\"}"
```

Ответ (пример):
```json
{
  "id": "f3b2c4e8d9a0410ba2c01fcd12345678",
  "status": "PENDING"
}
```

### 2. Проверка статуса и результата

```bash
curl "http://localhost:8080/api/v1/status/f3b2c4e8d9a0410ba2c01fcd12345678"
```

Пример ответа:
```json
{
  "id": "f3b2c4e8d9a0410ba2c01fcd12345678",
  "text": "Hello world. This is a test!",
  "status": "DONE",
  "created_at": "2026-02-17T10:00:00Z",
  "updated_at": "2026-02-17T10:00:01Z",
  "result": {
    "words_count": 5,
    "characters_count": 30,
    "sentences_count": 2,
    "average_word_length": 4.2,
    "original_request_id": "f3b2c4e8d9a0410ba2c01fcd12345678"
  }
}
```

### 3. Health‑checks

```bash
curl "http://localhost:8080/api/v1/health"
curl "http://localhost:8081/api/v1/health"
```

Оба сервиса возвращают:
```json
{"status":"ok"}
```

## Используемые технологии

- **Go 1.22** — реализация микросервисов.
- **net/http** — REST API и HTTP‑клиент.
- **Горутины** — асинхронный вызов сервиса B из сервиса A.
- **context** — управление таймаутами и graceful shutdown.
- **in‑memory хранилище** — `sync.RWMutex` + map для хранения запросов.
- **Docker, docker‑compose** — контейнеризация и запуск двух сервисов.

## Соответствие техническим требованиям

- HTTP‑взаимодействие между сервисами реализовано через `POST /api/v1/analyze`.
- Хранение данных — in‑memory репозиторий в сервисе A.
- Логирование — стандартный логгер `log.Logger` в обоих сервисах.
- Таймауты и ошибки сети — контекст с таймаутом при вызове сервиса B, обработка ошибок и логирование.
- Graceful shutdown — корректное завершение HTTP‑серверов по сигналам `SIGINT`/`SIGTERM`.
