# User Plants Service

Личная коллекция растений пользователя: CRUD, отметка полива, история поливов. Internal API для scheduler (растения к поливу и смена статуса).

**Base URL (локально):** `http://localhost:8082`

**Content-Type** для запросов с телом: `application/json`

---

## Аутентификация

Пользовательские эндпоинты требуют JWT от [user-service](../user-service/README.md):

```http
Authorization: Bearer <access_token>
```

`JWT_SECRET` должен совпадать с user-service.

Internal эндпоинты — заголовок:

```http
X-Internal-Key: <INTERNAL_API_KEY>
```

### Статусы растения

| Значение | Смысл |
|---|---|
| `ok` | полив не просрочен / после полива |
| `due` | пора полить (ставит scheduler) |
| `overdue` | полив просрочен (ставит scheduler) |

`nextWateringAt` = `lastWateredAt` + `wateringIntervalDays`. В needing-water попадают растения с `next_watering_at <= now`.

---

## Эндпоинты (пользователь)

### POST `/user-plants`

Создать растение в коллекции.

**Auth:** JWT

**Request body:**

```json
{
  "speciesId": 715,
  "name": "Grigory",
  "wateringIntervalDays": 7,
  "lastWateredAt": "2026-09-20T10:00:00Z"
}
```

`lastWateredAt` не может быть в будущем. Начальный `status` — `ok`.

**Response `201 Created`:**

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "speciesId": 715,
  "userId": "bab45d6d-eb91-4301-a081-5845d712bd81",
  "name": "Grigory",
  "wateringIntervalDays": 7,
  "status": "ok",
  "lastWateredAt": "2026-09-20T10:00:00Z",
  "nextWateringAt": "2026-09-27T10:00:00Z",
  "createdAt": "2026-10-03T05:00:00Z",
  "updatedAt": "2026-10-03T05:00:00Z"
}
```

**Errors:** `401` — нет/невалидный JWT; `400` — невалидные поля.

---

### GET `/user-plants`

Список растений текущего пользователя.

**Auth:** JWT

**Response `200 OK`:**

```json
{
  "userPlantsList": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "speciesId": 715,
      "userId": "bab45d6d-eb91-4301-a081-5845d712bd81",
      "name": "Grigory",
      "wateringIntervalDays": 7,
      "status": "ok",
      "lastWateredAt": "2026-09-20T10:00:00Z",
      "nextWateringAt": "2026-09-27T10:00:00Z",
      "createdAt": "2026-10-03T05:00:00Z",
      "updatedAt": "2026-10-03T05:00:00Z"
    }
  ]
}
```

---

### GET `/user-plants/{id}`

Одно растение по ID (только своё).

**Auth:** JWT

**Errors:** `404` — не найдено или чужое; `400` — невалидный UUID.

---

### PATCH `/user-plants/{id}`

Частичное обновление. Все поля опциональны.

**Auth:** JWT

**Request body (пример):**

```json
{
  "name": "New name",
  "wateringIntervalDays": 5,
  "status": "ok",
  "nextWateringAt": "2026-10-10T10:00:00Z"
}
```

**Response:** `204 No Content`

**Errors:** `404` — не найдено; `400` — невалидный status и т.п.

---

### DELETE `/user-plants/{id}`

Удалить растение.

**Auth:** JWT

**Response:** `204 No Content`

---

### POST `/user-plants/{id}/water`

Отметить полив: обновляет `lastWateredAt`, `nextWateringAt`, `status` → `ok`, пишет событие в историю.

**Auth:** JWT

**Request body (опционально):**

```json
{
  "wateredAt": "2026-10-03T12:00:00Z"
}
```

Если тело пустое — берётся текущее время.

**Response:** `204 No Content`

---

### GET `/user-plants/{id}/watering-events`

История поливов растения.

**Auth:** JWT

**Response `200 OK`:**

```json
[
  {
    "id": "...",
    "userPlantId": "550e8400-e29b-41d4-a716-446655440000",
    "wateredAt": "2026-10-03T12:00:00Z",
    "createdAt": "2026-10-03T12:00:01Z"
  }
]
```

---

## Эндпоинты (internal)

Для [scheduler-service](../scheduler-service/README.md). Не для фронта.

### GET `/internal/user-plants/needing-water`

Растения с `next_watering_at <= now`.

**Auth:** `X-Internal-Key`

**Response `200 OK`:** массив объектов `UserPlant` (как в create).

**Errors:** `401` — нет/неверный ключ.

---

### PATCH `/internal/user-plants/{id}/status`

Сменить статус (`ok` | `due` | `overdue`).

**Auth:** `X-Internal-Key`

**Request body:**

```json
{
  "status": "due"
}
```

**Response:** `204 No Content`

---

## Формат ошибок

Ошибки возвращаются как **plain text**, не JSON:

```http
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8

invalid interval
```

---

## Запуск

### Docker (из корня репозитория)

```bash
docker compose up -d --build
```

Нужны Postgres и совпадающие `JWT_SECRET` / `INTERNAL_API_KEY` в корневом `.env`.

### Локально (только Go)

1. Поднять Postgres (`docker compose up -d postgres`).
2. Скопировать `.env.example` → `.env`.
3. Запустить:

```bash
go run ./cmd/server
```

---

## CORS

CORS на бэкенде пока **не настроен**. Если фронт на другом origin, нужен proxy или middleware CORS.

---

## Переменные окружения

См. `.env.example`:

| Переменная | Описание |
|---|---|
| `ADDR` | Адрес сервера (по умолчанию `:8082`) |
| `DATABASE_URL` | PostgreSQL connection string |
| `JWT_SECRET` | Секрет JWT (**обязателен**, как в user-service) |
| `INTERNAL_API_KEY` | Ключ для `/internal/*` (**обязателен**) |
