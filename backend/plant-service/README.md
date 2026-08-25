# Plant Service

Каталог растений: список видов и детальная карточка. Данные из [Perenual](https://perenual.com/); в details дополнительно подтягивается описание из Wikipedia (по scientific name).

**Base URL (локально):** `http://localhost:8081`

Каталог **публичный** — JWT не требуется.

**Content-Type** ответов: `application/json; charset=utf-8`

---

## Источники данных

| Эндпоинт | Perenual | Wikipedia |
|---|---|---|
| `GET /plants` | да (species-list) | нет |
| `GET /plants/{id}` | да (species details) | да, если есть scientific name |

---

## Эндпоинты

### GET `/plants`

Список видов с пагинацией и опциональным поиском.

**Auth:** не требуется

**Query params:**

| Параметр | Описание |
|---|---|
| `page` | Номер страницы (по умолчанию `1`; меньше 1 приводится к `1`) |
| `query` | Поисковая строка (передаётся в Perenual) |

**Пример:** `GET /plants?page=1&query=monstera`

**Response `200 OK`:**

```json
{
  "species": [
    {
      "id": 715,
      "common_name": "Swiss Cheese Plant",
      "scientific_names": ["Monstera deliciosa"],
      "other_names": ["Split Leaf Philodendron"],
      "image": {
        "regular_url": "https://...",
        "small_url": "https://..."
      }
    }
  ],
  "per_page": 30,
  "current_page": 1,
  "last_page": 10,
  "total": 300
}
```

Поле `image` может быть `null`, если у вида нет картинки.

**Errors:**

| Код | Когда |
|---|---|
| `400` | `page` не число |
| `500` | ошибка внешнего API / сети |

---

### GET `/plants/{id}`

Детальная карточка вида по ID Perenual.

**Auth:** не требуется

**Path params:** `id` — целое число (ID вида)

**Пример:** `GET /plants/715`

**Response `200 OK`:** (поля могут быть пустыми / `null`, если Perenual/Wikipedia их не вернули)

```json
{
  "id": 715,
  "common_name": "Swiss Cheese Plant",
  "scientific_names": ["Monstera deliciosa"],
  "other_names": ["Split Leaf Philodendron"],
  "image": {
    "regular_url": "https://...",
    "small_url": "https://..."
  },
  "growth_rate": "High",
  "type": "vine",
  "dimensions": {
    "min_value": 1,
    "max_value": 3,
    "unit": "meters"
  },
  "cycle": "Perennial",
  "watering": "Average",
  "watering_general_benchmark": {
    "value": "7-10",
    "unit": "days"
  },
  "sunlight": ["part shade", "filtered shade"],
  "pruning_months": ["March", "April"],
  "pruning_count": {
    "amount": 1,
    "interval": "yearly"
  },
  "soil": ["Well-draining"],
  "maintenance": "Moderate",
  "poisonous_to_humans": true,
  "poisonous_to_pets": true,
  "care_level": "Moderate",
  "perenual_description": "...",
  "flowers": true,
  "flowering_season": "Spring",
  "wikipedia_description": "...",
  "wikipedia_extract": "...",
  "wikipedia_extract_html": "<p>...</p>",
  "wikipedia_url": "https://en.wikipedia.org/wiki/Monstera_deliciosa"
}
```

Поля Wikipedia (`wikipedia_*`) заполняются только если удалось найти страницу по scientific name; иначе `null`.

**Errors:**

| Код | Когда |
|---|---|
| `400` | `id` не число |
| `500` | ошибка внешнего API / сети |

---

## Формат ошибок

Ошибки возвращаются как **plain text**, не JSON:

```http
HTTP/1.1 400 Bad Request
Content-Type: text/plain; charset=utf-8

strconv.Atoi: parsing "abc": invalid syntax
```

На фронте читай `response.text()` или ориентируйся на HTTP-код.

---

## Запуск

### Docker (из корня репозитория)

В корневом `.env` (рядом с `docker-compose.yaml`) задай ключ:

```env
PERENUAL_KEY=your_key
```

Затем:

```bash
docker compose up -d --build
```

Сервис доступен на `http://localhost:8081`.

### Локально (только Go)

1. Скопировать `.env.example` → `.env` и задать `PERENUAL_KEY`.
2. Запустить сервис:

```bash
go run ./cmd/server
```

Postgres не нужен — plant-service ходит только во внешние API.

---

## CORS

CORS на бэкенде пока **не настроен**. Если фронт на другом origin (например `http://localhost:5173`), понадобится proxy в dev-сервере фронта или middleware CORS на plant-service.

---

## Переменные окружения

См. `.env.example`:

| Переменная | Описание |
|---|---|
| `ADDR` | Адрес сервера (по умолчанию `:8080`; в Docker обычно `:8081`) |
| `TIMEOUT` | Таймаут HTTP-клиентов к Perenual/Wikipedia (например `10s`; по умолчанию `10s`) |
| `PERENUAL_KEY` | API-ключ Perenual (**обязателен**) |
