# pLaNtS

Сервис для ухода за растениями.

## Сервисы

### [User Service](backend/user-service/README.md)

Микросервис аутентификации на Go: регистрация, вход, JWT-токены, профиль пользователя.  
Порт по умолчанию: `8080`.

## Запуск

Из корня репозитория:

```bash
docker compose up -d --build
```

Поднимает Postgres и user-service. Подробности — в README каждого сервиса.
