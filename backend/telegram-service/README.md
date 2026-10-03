# Telegram Service

Привязка аккаунта pLaNtS к Telegram-боту и отправка уведомлений. Пользователь получает deep link, жмёт Start в боте — сервис сохраняет `chat_id`. Уведомления приходят по HTTP (internal) или из Kafka (топик `watering.notify`).

**Base URL (локально):** `http://localhost:8084`

**Content-Type** для запросов с телом: `application/json`

**Health:** `GET /health` → `ok`

---

## Flow привязки

1. Клиент с JWT вызывает `POST /telegram/link` → получает `linkUrl`.
2. Пользователь открывает ссылку в Telegram и жмёт **Start**.
3. Бот (long polling) ловит `/start <code>`, сохраняет `chat_id` за пользователем.
4. После этого работают `POST /notify` и consumer Kafka.

---

## Аутентификация

| Эндпоинт | Auth |
|---|---|
| `POST /telegram/link` | `Authorization: Bearer <access_token>` (JWT, тот же `JWT_SECRET`) |
| `POST /notify` | `X-Internal-Key: <INTERNAL_API_KEY>` |
| Kafka consumer | без HTTP; читает топик внутри процесса |

---

## Эндпоинты

### POST `/telegram/link`

Создать (или обновить) код привязки и URL для бота.

**Auth:** JWT

**Request body:** нет

**Response `200 OK`:**

```json
{
  "linkUrl": "https://t.me/YourBotName?start=abc123...",
  "linkCode": "abc123..."
}
```

Открой `linkUrl` в Telegram и нажми Start.

**Errors:** `401` — нет/невалидный JWT; `500` — ошибка БД и т.п.

---

### POST `/notify`

Отправить текст пользователю в Telegram (если chat привязан). Удобно для ручных проверок; основной путь уведомлений о поливе — Kafka.

**Auth:** `X-Internal-Key`

**Request body:**

```json
{
  "userId": "bab45d6d-eb91-4301-a081-5845d712bd81",
  "text": "Пора полить растение «Grigory»"
}
```

**Response:** `204 No Content`

**Errors:**

| Код | Когда |
|---|---|
| `401` | нет/неверный `X-Internal-Key` |
| `400` | невалидный JSON / `userId` / пустой text |
| `404` | Telegram не привязан (`telegram is not linked`) |

---

### GET `/health`

**Auth:** не требуется

**Response `200 OK`:** тело `ok` (plain text).

---

## Kafka consumer

При старте сервиса поднимается consumer group (по умолчанию `telegram-notifier`) на топик `watering.notify`.

Сообщение (JSON), которое публикует [scheduler-service](../scheduler-service/README.md):

```json
{
  "user_id": "bab45d6d-eb91-4301-a081-5845d712bd81",
  "plant_id": "550e8400-e29b-41d4-a716-446655440000",
  "plant_name": "Grigory",
  "status": "due",
  "text": "Пора полить растение «Grigory»",
  "occurred_at": "2026-10-03T12:00:00Z"
}
```

Consumer вызывает тот же `Notify(user_id, text)`. Битый JSON коммитится и пропускается; если пользователь не привязан — сообщение тоже коммитится (без бесконечного ретрая).

Проверка вручную:

```bash
docker exec -i plants-redpanda rpk topic produce watering.notify <<'EOF'
{"user_id":"<uuid>","plant_id":"<uuid>","plant_name":"Test","status":"due","text":"Проверка","occurred_at":"2026-10-03T12:00:00Z"}
EOF
```

---

## Формат ошибок

Ошибки HTTP — **plain text**, не JSON:

```http
HTTP/1.1 404 Not Found
Content-Type: text/plain; charset=utf-8

telegram is not linked
```

---

## Запуск

### Docker (из корня репозитория)

В корневом `.env`:

```env
TELEGRAM_BOT_TOKEN=...
TELEGRAM_BOT_USERNAME=YourBotName
JWT_SECRET=...
INTERNAL_API_KEY=change-me-internal
```

```bash
docker compose up -d --build
```

Нужны Postgres и Redpanda. Из контейнера brokers: `redpanda:9092`.

### Локально (только Go)

1. Postgres + Redpanda (`docker compose up -d postgres redpanda`).
2. Скопировать `.env.example` → `.env`. Для локального Go к Redpanda с хоста обычно `KAFKA_BROKERS=localhost:19092`.
3. Запустить:

```bash
go run ./cmd/server
```

---

## CORS

CORS на бэкенде пока **не настроен**.

---

## Переменные окружения

См. `.env.example`:

| Переменная | Описание |
|---|---|
| `ADDR` | Адрес сервера (по умолчанию `:8084`) |
| `DATABASE_URL` | PostgreSQL connection string |
| `TELEGRAM_BOT_TOKEN` | Токен бота от BotFather (**обязателен**) |
| `TELEGRAM_BOT_USERNAME` | Username бота без `@` или с ним (**обязателен**) |
| `JWT_SECRET` | Секрет JWT (**обязателен**) |
| `INTERNAL_API_KEY` | Ключ для `/notify` (**обязателен**) |
| `TIMEOUT` | Таймаут HTTP к Telegram API (например `40s`) |
| `KAFKA_BROKERS` | Адреса брокера (по умолчанию `redpanda:9092`) |
| `KAFKA_TOPIC_WATERING` | Топик событий (по умолчанию `watering.notify`) |
| `KAFKA_GROUP_ID` | Consumer group (по умолчанию `telegram-notifier`) |
