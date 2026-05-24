# Repair CRM

SaaS-платформа для сервисных центров по ремонту телефонов: публичные заявки без обязательной регистрации, постоянный трекинг в БД, фото приёмки, цифровые гарантийные и скидочные талоны, реферальная программа, интеграция API и тарифы подписки.

## Возможности

| Область | Описание |
|--------|----------|
| **Клиент** | Заявка без пароля; телефон необязателен; IMEI обязателен (15 цифр); личная ссылка `/order/{token}` сохраняется в браузере и в PostgreSQL |
| **Трекинг** | Статусы от «Заявка» до «Выдано»; после принятия — публичная ссылка `/track/{hash}` |
| **Безопасность** | IMEI в публичном API маскируется (•••••••••••1234); фото приёмки загружается один раз и не редактируется |
| **Отзывы** | Только клиент, только после статуса «Выдано»; мастер видит отзыв, но не может его отправить |
| **Талоны** | При выдаче генерируются уникальные коды гарантии (`W…`) и скидки (`D…`) — остаются у клиента на странице заказа |
| **Рефералы** | У каждой заявки код; при выдаче рефералу начисляется дополнительный скидочный талон |
| **Админ** | Панель `/admin`, настройки мастерской, смена API-ключа, демо-подписка на тариф |
| **SaaS** | Регистрация мастерской на `/pricing`, тарифы `starter` / `pro` / `business` |
| **Интеграция** | `POST /api/integration/tickets` с заголовком `X-API-Key` |

## Стек

- **Backend:** Go 1.21, Chi, PostgreSQL 16, JWT, bcrypt
- **Frontend:** React 18, TypeScript, Vite, Tailwind CSS
- **Инфра:** Docker Compose (postgres, backend, nginx)

## Требования

- Docker и Docker Compose **или**
- Go 1.21+, Node.js 20+, PostgreSQL 16

## Быстрый старт (Docker)

### 1. Клонирование и переход в каталог

```bash
cd /home/almas/repair-crm-3
```

### 2. Переменные окружения (опционально)

Скопируйте пример и при необходимости измените секреты:

```bash
cp .env.example .env
```

Для production обязательно смените `JWT_SECRET` и API-ключ мастерской.

### 3. Миграции и демо-данные

```bash
make seed-demo
```

`seed-demo` сам поднимает PostgreSQL, применяет миграции (с учётом уже выполненных) и создаёт демо-мастерскую.

> **Важно:** после `make clean` нельзя сразу только `docker compose up` — без миграций таблиц не будет. Используйте `make seed-demo` или `make run`.

`seed-demo` создаёт:

- мастерскую `Demo Repair` (id=1);
- мастера `admin` / `admin123`;
- демо API-ключ (см. миграцию `000003`).

### 4. Запуск всех сервисов

```bash
docker compose up --build
```

или одной командой:

```bash
make run
```

### 5. Открыть в браузере

| URL | Назначение |
|-----|------------|
| http://localhost:3000 | Публичная заявка |
| http://localhost:3000/my-orders | Сохранённые заявки клиента |
| http://localhost:3000/pricing | Тарифы и регистрация мастерской |
| http://localhost:3000/admin | Панель мастера |
| http://localhost:8080/health | Healthcheck API |

**Демо-вход:** `admin` / `admin123`

## Локальная разработка без Docker

### PostgreSQL

```bash
# пример: база repair_crm, пользователь repair / repair
export DATABASE_URL="postgres://repair:repair@localhost:5432/repair_crm?sslmode=disable"
```

Примените миграции вручную:

```bash
for f in migrations/*.up.sql; do psql "$DATABASE_URL" -v ON_ERROR_STOP=1 -f "$f"; done
```

Затем `make seed-demo` (через docker postgres) или выполните SQL из `Makefile`.

### Backend

```bash
export JWT_SECRET="dev-secret-change-me"
export DEFAULT_WORKSHOP_ID=1
export CORS_ALLOWED_ORIGIN="*"
go run ./cmd/app
```

Слушает `:8080`.

### Frontend

```bash
cd frontend
npm ci
export VITE_API_URL="http://localhost:8080/api"
npm run dev
```

Vite по умолчанию: http://localhost:5173 — настройте `CORS_ALLOWED_ORIGIN` на этот адрес.

## Сценарий «клиент → мастер → выдача»

1. Клиент на `/` заполняет имя, IMEI, модель, поломку (телефон по желанию).
2. После отправки получает ссылку `/order/{client_token}` — она сохраняется в `localStorage` и в БД.
3. Мастер в `/admin` видит заявку «Заявка», загружает **фото приёмки** (камера или файл), переводит в «Принято».
4. Клиент на странице заказа видит обновление статуса (можно обновить страницу).
5. Мастер ведёт ремонт до «Выдано» — клиенту выдаются коды гарантии и скидки.
6. Клиент оставляет **отзыв** (1–5 звёзд) — только он, только после «Выдано».
7. Если друг указал реферальный код при заявке — после выдачи рефереру начисляется бонусный скидочный код.

## IMEI и регистрация

- **Регистрация клиенту не нужна.** Идентификация: секретный `client_token` + IMEI в базе.
- **IMEI обязателен** при подаче заявки — стандарт для сервисных центров (учёт устройства, споры, мошенничество).
- В публичном трекинге IMEI **не показывается целиком** — только маска с последними 4 цифрами.
- Полный IMEI доступен только мастеру в админ-панели.

## API (основное)

### Публичные

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/api/applications` | Новая заявка → `{ ticket, order_url }` |
| GET | `/api/orders/{token}` | Статус заявки по client_token |
| POST | `/api/orders/{token}/review` | Отзыв клиента |
| GET | `/api/track/{hash}` | Трекинг по short_hash |
| POST | `/api/track/{hash}/review` | Отзыв по hash |
| GET | `/api/plans` | Список тарифов |
| POST | `/api/workshops/register` | Регистрация мастерской + JWT |

### Интеграция

```bash
curl -X POST http://localhost:8080/api/integration/tickets \
  -H "Content-Type: application/json" \
  -H "X-API-Key: demo-integration-key-change-me" \
  -d '{
    "client_name": "Иван",
    "client_phone": "+77001234567",
    "imei": "123456789012345",
    "brand": "Samsung",
    "model": "Galaxy S21",
    "defect_description": "Не включается"
  }'
```

### Мастер (JWT `Authorization: Bearer …`)

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/api/auth/login` | Вход |
| GET/PUT | `/api/workshop/settings` | Настройки мастерской |
| POST | `/api/workshop/subscribe` | Активация тарифа (демо) |
| POST | `/api/workshop/integration-key/rotate` | Новый API-ключ |
| GET/POST/PUT/DELETE | `/api/tickets…` | CRUD заявок |
| POST | `/api/tickets/{id}/intake-photo` | Фото приёмки (один раз) |
| GET | `/api/tickets/{id}/intake-photo` | Просмотр фото (только JWT) |

## Переменные окружения backend

| Переменная | По умолчанию | Описание |
|------------|--------------|----------|
| `HTTP_ADDR` | `:8080` | Адрес HTTP |
| `DATABASE_URL` | `postgres://repair:repair@postgres:5432/repair_crm?sslmode=disable` | PostgreSQL |
| `JWT_SECRET` | `dev-secret-change-me` | Секрет JWT |
| `CORS_ALLOWED_ORIGIN` | `*` | CORS |
| `DEFAULT_WORKSHOP_ID` | `1` | Мастерская для публичных заявок с `/` |

## Команды Make

```bash
make migrate        # применить миграции (идемпотентно)
make seed-demo      # миграции + Demo Repair + admin/admin123
make run            # migrate + docker compose up --build
make e2e-test       # smoke-тест API (backend должен быть запущен)
make backend-test   # go test ./...
make frontend-build # production-сборка SPA
make down           # остановить контейнеры
make clean          # удалить контейнеры и volume БД
```

## Проверка

```bash
make e2e-test
```

```bash
curl http://localhost:8080/health
# {"status":"ok"}

curl -X POST http://localhost:8080/api/applications \
  -H 'Content-Type: application/json' \
  -d '{
    "client_name":"Алина",
    "client_phone":"",
    "imei":"356938035643809",
    "brand":"Apple",
    "model":"iPhone 13",
    "defect_description":"Разбит экран"
  }'
```

В ответе будут `ticket.client_token` и `order_url`.

## Миграции

| Файл | Содержание |
|------|------------|
| `000001_init` | Базовые таблицы |
| `000002_public_flow` | Русские статусы, поля клиента на заявке |
| `000003_production_features` | client_token, фото, талоны, рефералы, тарифы, настройки мастерской |

## Память (Docker)

- backend: `512m`, `GOMEMLIMIT=460MiB`
- PostgreSQL: `512m`
- frontend nginx: `256m`

## Продакшен (чеклист)

1. Сменить `JWT_SECRET`, `integration_api_key`, пароли БД.
2. Подключить реальный платёжный шлюз вместо демо `POST /api/workshop/subscribe`.
3. Настроить HTTPS и `CORS_ALLOWED_ORIGIN` на домен фронтенда.
4. Резервное копирование PostgreSQL (фото приёмки хранятся в `BYTEA`).
5. Ограничить размер загрузок на reverse proxy (фото ≤ 2 МБ на уровне API).

## Структура проекта

```
cmd/app/              # точка входа HTTP-сервера
internal/api/         # handlers, router, middleware
internal/models/      # доменные типы
internal/repository/  # SQL
internal/service/     # бизнес-логика
pkg/                  # jwt, password, hash, codes
migrations/           # SQL миграции
frontend/src/         # React SPA
```

## Лицензия

Проприетарный MVP — уточните у владельца репозитория.
