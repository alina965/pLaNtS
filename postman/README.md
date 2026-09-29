# Postman

## Import

1. Open Postman → **Import**
2. Add:
   - `pLaNtS.postman_collection.json`
   - `pLaNtS.local.postman_environment.json`
3. Select environment **pLaNtS Local** (top-right)

## Suggested order

1. **1. Auth** → Register or Login (saves `access_token`)
2. **Get Me** (saves `user_id`)
3. **2. Plants catalog**
4. **3. User plants** → Create (saves `plant_id`)
5. **5. Telegram** → Create link → open `telegram_link_url` in Telegram → Start
6. **Notify user**
7. **4. Internal** → needing-water / update status

## Variables

| Variable | Meaning |
|----------|---------|
| `email` / `password` | change before Register if email already used |
| `last_watered_at` | collection variable, use a past ISO date so plant can become due |
| `internal_api_key` | must match compose / `.env` (`change-me-internal` by default) |
| `species_id` | Perenual id, default `715` |

Scheduler cron is not triggered from Postman — check logs / wait for 09:00 or temporarily set cron to `*/1 * * * *`.
