# pLaNtS

Сервис для ухода за растениями.

## Сервисы

### [User Service](backend/user-service/README.md)

Микросервис аутентификации на Go: регистрация, вход, JWT-токены, профиль пользователя.  
Порт по умолчанию: `8080`.

### [Plant Service](backend/plant-service/README.md)

Каталог растений на Go: список и детали видов (Perenual + Wikipedia).  
Порт по умолчанию: `8081`.

## Запуск

Из корня репозитория задай секреты в `.env` (см. сервисы), затем:

```bash
docker compose up -d --build
```

Поднимает Postgres, user-service и plant-service. Подробности — в README каждого сервиса.
