# pLaNtS

Сервис для ухода за растениями: каталог видов, личная коллекция, напоминания о поливе через Telegram.

## Сервисы

| Сервис | Порт | Описание |
|--------|------|----------|
| [user-service](backend/user-service/README.md) | `8080` | Регистрация, вход, JWT, профиль |
| [plant-service](backend/plant-service/README.md) | `8081` | Каталог видов (Perenual + Wikipedia) |
| user-plants-service | `8082` | Растения пользователя, полив, internal API для scheduler |
| scheduler-service | `8083` | Cron: находит растения к поливу и шлёт уведомления |
| telegram-service | `8084` | Привязка аккаунта к боту и отправка сообщений |
| postgres | `5432` | Общая БД |

## Запуск

В корне репозитория создай `.env` (можно от `.env.example`) и задай:

```env
PERENUAL_KEY=...
JWT_SECRET=...
INTERNAL_API_KEY=change-me-internal
TELEGRAM_BOT_TOKEN=...
TELEGRAM_BOT_USERNAME=...
```

Опционально: `CRON_LOCATION` (по умолчанию `Europe/Moscow`).

```bash
docker compose up -d --build
```

Поднимает Postgres и все backend-сервисы. Миграции накатываются при первом создании volume Postgres.

Проверка scheduler / telegram:

```bash
curl http://localhost:8083/health
curl http://localhost:8084/health
```

## Как связаны сервисы

1. Пользователь регистрируется через **user-service** и получает JWT.
2. Каталог смотрит в **plant-service**, свои растения ведёт в **user-plants-service**.
3. Через **telegram-service** пользователь привязывает Telegram (deep link `/start`).
4. **scheduler-service** по cron запрашивает растения к поливу у user-plants и шлёт уведомления через telegram-service (`X-Internal-Key`).
