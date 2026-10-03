# Scheduler Service

Фоновый сервис: по cron находит растения к поливу в [user-plants-service](../user-plants-service/README.md), обновляет статус и публикует событие в Kafka. [telegram-service](../telegram-service/README.md) читает топик и шлёт сообщение в бот.

**Base URL (локально):** `http://localhost:8083` (только health)

**Health:** `GET /health` → `ok`

Публичного REST API для полива нет — работа идёт по расписанию.

---

## Как работает

1. Cron (по умолчанию каждую минуту: `* * * * *`, таймзона `CRON_LOCATION`).
2. `GET /internal/user-plants/needing-water` у user-plants (`X-Internal-Key`).
3. Для каждого растения:
   - `ok` → статус `due`, текст «Пора полить…»
   - `due` → статус `overdue`, текст «Полив просрочен…»
   - иначе — пропуск
4. `PATCH /internal/user-plants/{id}/status`
5. Produce в топик `watering.notify` (не ждёт Telegram).

```text
scheduler → user-plants (HTTP internal)
         → Redpanda / Kafka (produce)
telegram-service ← consume ← топик
```

---

## Эндпоинты

### GET `/health`

**Auth:** не требуется

**Response `200 OK`:** тело `ok` (plain text).

---

## Событие Kafka

Топик по умолчанию: `watering.notify`

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

Просмотр сообщений:

```bash
docker exec -it plants-redpanda rpk topic consume watering.notify -f '%v\n'
```

---

## Проверка E2E (кратко)

1. Создай растение с `lastWateredAt` в прошлом так, чтобы `nextWateringAt` уже наступил (см. user-plants README).
2. Привяжи Telegram (`POST /telegram/link` → Start в боте).
3. Дождись тика cron (до ~1 минуты при `* * * * *`).
4. В боте — текст про полив; статус растения — `due` или `overdue`.

Логи:

```bash
docker logs -f plants-scheduler-service
docker logs -f plants-telegram-service
```

---

## Запуск

### Docker (из корня репозитория)

```bash
docker compose up -d --build
```

Нужны user-plants-service, Redpanda (`KAFKA_BROKERS=redpanda:9092`), общий `INTERNAL_API_KEY`.

### Локально (только Go)

1. Поднять зависимости (`postgres`, `user-plants-service`, `redpanda`).
2. Скопировать `.env.example` → `.env`. С хоста: `USER_PLANTS_URL=http://localhost:8082`, `KAFKA_BROKERS=localhost:19092`.
3. Запустить:

```bash
go run ./cmd/server
```

`TELEGRAM_URL` в config пока **обязателен при старте** (compose его прокидывает), но в коде notify через HTTP больше не вызывается — только Kafka.

---

## Переменные окружения

См. `.env.example`:

| Переменная | Описание |
|---|---|
| `ADDR` | Адрес health-сервера (по умолчанию `:8083`) |
| `USER_PLANTS_URL` | Base URL user-plants (**обязателен**) |
| `TELEGRAM_URL` | Пока обязателен в config; для notify не используется (остался от HTTP-клиента) |
| `TIMEOUT` | Таймаут HTTP к user-plants в секундах (по умолчанию `30`) |
| `INTERNAL_API_KEY` | Ключ internal API (**обязателен**) |
| `CRON_LOCATION` | Таймзона cron (по умолчанию `Europe/Moscow`) |
| `KAFKA_BROKERS` | Брокер (по умолчанию `redpanda:9092`) |
| `KAFKA_TOPIC_WATERING` | Топик (по умолчанию `watering.notify`) |
