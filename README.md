# pLaNtS

Сервис для ухода за растениями: каталог видов, личная коллекция, напоминания о поливе через Telegram.

## Сервисы

| Сервис | Порт | Описание |
|--------|------|----------|
| [user-service](backend/user-service/README.md) | `8080` | Регистрация, вход, JWT, профиль |
| [plant-service](backend/plant-service/README.md) | `8081` | Каталог видов (Perenual + Wikipedia), кэш в Redis |
| [user-plants-service](backend/user-plants-service/README.md) | `8082` | Растения пользователя, полив, internal API |
| [scheduler-service](backend/scheduler-service/README.md) | `8083` | Cron: needing-water → Kafka |
| [telegram-service](backend/telegram-service/README.md) | `8084` | Привязка бота, notify (HTTP + Kafka consumer) |
| postgres | `5432` | Общая БД |
| redis | `6379` | Кэш ответов Perenual |
| redpanda | `19092` (host) / `9092` (docker) | Брокер Kafka API для уведомлений |

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

Поднимает Postgres, Redis, Redpanda и все backend-сервисы. Миграции накатываются при первом создании volume Postgres.

Проверка:

```bash
curl http://localhost:8083/health
curl http://localhost:8084/health
```

Если Postgres volume уже существовал до новых миграций — таблицы telegram/user-plants могут отсутствовать; тогда нужен новый volume или ручной SQL из `backend/*/migrations`.

## Как связаны сервисы

1. Пользователь регистрируется через **user-service** и получает JWT.
2. Каталог смотрит в **plant-service** (Perenual + Redis-кэш), свои растения ведёт в **user-plants-service**.
3. Через **telegram-service** пользователь привязывает Telegram (deep link `/start`).
4. **scheduler-service** по cron запрашивает растения к поливу у user-plants и публикует событие в Kafka (`watering.notify`).
5. **telegram-service** читает топик и шлёт сообщение в бот.
