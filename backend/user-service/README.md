# User Service

Сервис аутентификации: регистрация, вход, обновление токенов, профиль текущего пользователя.

**Base URL (локально):** `http://localhost:8080`

**Content-Type** для запросов с телом: `application/json`

---

## Аутентификация

После регистрации или логина клиент получает пару токенов:

| Поле | Описание |
|---|---|
| `access_token` | JWT для защищённых запросов. Короткий срок жизни. |
| `refresh_token` | UUID для обновления access без пароля. Длинный срок. |
| `expires_at` | Когда истечёт `access_token` (ISO 8601). |

**Защищённые эндпоинты** требуют заголовок:

```http
Authorization: Bearer <access_token>
```

### Типичный flow на фронте

1. `POST /auth/login` или `POST /auth/register` → сохранить оба токена.
2. Все запросы к API → `Authorization: Bearer <access_token>`.
3. Если ответ `401` и access истёк → `POST /auth/refresh` → сохранить новую пару → повторить запрос.
4. `GET /me` → данные текущего пользователя.

---

## Эндпоинты

### POST `/auth/register`

Регистрация нового пользователя. Сразу возвращает токены (как после логина).

**Auth:** не требуется

**Request body:**

```json
{
  "email": "user@example.com",
  "password": "qwerty123",
  "first_name": "Иван",
  "last_name": "Иванов"
}
```

**Response `201 Created`:**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "550e8400-e29b-41d4-a716-446655440000",
  "expires_at": "2026-08-23T17:15:03+07:00"
}
```

**Errors:** `400` — невалидный JSON, email уже занят и т.п. (текст ошибки в теле).

---

### POST `/auth/login`

Вход по email и паролю.

**Auth:** не требуется

**Request body:**

```json
{
  "email": "user@example.com",
  "password": "qwerty123"
}
```

**Response `200 OK`:**

```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIs...",
  "refresh_token": "550e8400-e29b-41d4-a716-446655440000",
  "expires_at": "2026-08-23T17:15:03+07:00"
}
```

**Errors:** `400` — неверный email или пароль (текст ошибки в теле).

---

### POST `/auth/refresh`

Выдаёт новую пару токенов по `refresh_token`. Старый refresh после успешного запроса больше не действует.

**Auth:** не требуется (передаётся `refresh_token` в теле)

**Request body:**

```json
{
  "refresh_token": "550e8400-e29b-41d4-a716-446655440000"
}
```

**Response `200 OK`:** тот же формат, что у login/register (`access_token`, `refresh_token`, `expires_at`).

**Errors:** `400` — токен не найден, отозван или истёк.

---

### GET `/me`

Профиль пользователя, которому принадлежит access token.

**Auth:** обязателен — `Authorization: Bearer <access_token>`

**Request body:** нет

**Response `200 OK`:**

```json
{
  "id": "bab45d6d-eb91-4301-a081-5845d712bd81",
  "email": "user@example.com",
  "first_name": "Иван",
  "last_name": "Иванов"
}
```

**Errors:**

| Код | Когда |
|---|---|
| `401` | нет заголовка, неверный или просроченный access token |
| `400` | пользователь не найден в БД |

---

## Формат ошибок

Сейчас ошибки возвращаются как **plain text**, не JSON:

```http
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8

wrong password
```

На фронте читай `response.text()` для сообщения или ориентируйся только на HTTP-код.

---

## Запуск

### Docker (из корня репозитория)

```bash
docker compose up -d --build
```

### Локально (только Go)

1. Поднять Postgres (`docker compose up -d postgres` из корня или свой инстанс).
2. Скопировать `.env.example` → `.env` и задать `JWT_SECRET`.
3. Запустить сервис:

```bash
go run ./cmd/server
```

---

## CORS

CORS на бэкенде пока **не настроен**. Если фронт на другом origin (например `http://localhost:5173`), понадобится proxy в dev-сервере фронта или middleware CORS на user-service.

---

## Переменные окружения

См. `.env.example`:

| Переменная | Описание |
|---|---|
| `ADDR` | Адрес сервера (по умолчанию `:8080`) |
| `DATABASE_URL` | PostgreSQL connection string |
| `JWT_SECRET` | Секрет для подписи JWT (**обязателен**) |
| `ACCESS_TOKEN_TTL` | Срок access token (например `15m`) |
| `REFRESH_TOKEN_TTL` | Срок refresh token (например `168h`) |
| `BCRYPT_COST` | Cost для bcrypt (10+) |
