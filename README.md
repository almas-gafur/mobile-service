# Repair CRM

SaaS MVP для мастерских по ремонту телефонов: публичные заявки без регистрации, цифровые гарантии, трекинг этапов ремонта и защищенные отзывы только после выдачи аппарата.

## Что внутри

- Backend: Go 1.21, `chi`, `database/sql`, `pgx`, JWT, bcrypt.
- DB: PostgreSQL 16, параметризованные SQL-запросы, транзакции для операций `device + ticket`.
- Frontend: React, TypeScript, Tailwind CSS.
- Infra: Docker Compose, лимиты памяти `512m` для backend и PostgreSQL.

## Быстрый запуск

1. Перейдите в проект:

```bash
cd /home/almas/project/repair-crm
```

2. Поднимите PostgreSQL и примените миграции:

```bash
make migrate
```

Команда запускает контейнер `postgres`, ждет готовности БД и последовательно применяет все файлы `migrations/*.up.sql`.

3. Создайте тестовую мастерскую и мастера:

```bash
make seed-demo
```

Эта команда выполняет SQL:

```sql
INSERT INTO workshops (id, name)
VALUES (1, 'Demo Repair')
ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name;

INSERT INTO masters (workshop_id, username, password_hash)
VALUES (1, 'admin', '$2b$12$xxHRJWHOPwTGxH.a/0xKz.2fWSAaA7t7tzwyZHDsG0j5dr4UCB6rK')
ON CONFLICT (username) DO UPDATE
SET workshop_id = EXCLUDED.workshop_id,
    password_hash = EXCLUDED.password_hash;
```

Демо-доступ: `admin / admin123`.

4. Запустите весь продукт:

```bash
docker compose up --build
```

Альтернативно одной командой для обычного запуска:

```bash
make run
```

## Где открыть

- Публичная заявка клиента: http://localhost:3000
- Панель мастера: http://localhost:3000/admin
- Backend healthcheck: http://localhost:8080/health
- PostgreSQL: `localhost:5432`, база `repair_crm`, пользователь `repair`, пароль `repair`

## Рабочий сценарий

1. Клиент открывает `/`, вводит имя, телефон, модель и описание поломки.
2. Backend создает заявку со статусом `Заявка` без пароля клиента.
3. Мастер входит в `/admin`, видит новую заявку и переводит ее в `Принято`.
4. При первом переводе из `Заявка` backend генерирует криптостойкий `short_hash` из 8 символов.
5. Мастер копирует ссылку `/track/{hash}` и отправляет клиенту.
6. Клиент видит прогресс ремонта, гарантию, дисклеймер по жидкости и промокод `FIX10`.
7. Отзыв можно отправить только когда мастер выставил статус `Выдано`; backend повторно проверяет статус перед записью.

## API

- `POST /api/applications` — публичная заявка.
- `POST /api/auth/login` — вход мастера.
- `GET /api/tickets` — список заявок мастерской, JWT.
- `POST /api/tickets` — ручное создание заявки мастером, JWT.
- `GET /api/tickets/{id}` — заявка, JWT.
- `PUT /api/tickets/{id}` — обновление заявки и статуса, JWT.
- `DELETE /api/tickets/{id}` — удаление заявки, JWT.
- `GET /api/track/{hash}` — публичный трекинг.
- `POST /api/track/{hash}/review` — отзыв клиента после статуса `Выдано`.

## Команды

```bash
make migrate        # применить миграции
make seed-demo      # создать Demo Repair и admin/admin123
make run            # миграции + запуск всех контейнеров
make backend-test   # go test ./...
make frontend-build # npm ci && npm run build
make down           # остановить контейнеры
make clean          # удалить контейнеры и volume с БД
```

## Проверка работоспособности

```bash
curl http://localhost:8080/health
```

Ожидаемый ответ:

```json
{"status":"ok"}
```

Публичную заявку можно проверить так:

```bash
curl -X POST http://localhost:8080/api/applications \
  -H 'Content-Type: application/json' \
  -d '{"client_name":"Алина","client_phone":"+77001234567","brand":"Apple","model":"iPhone 13","defect_description":"Разбит экран"}'
```

## Переменные окружения backend

- `HTTP_ADDR`, default `:8080`
- `DATABASE_URL`
- `JWT_SECRET`
- `CORS_ALLOWED_ORIGIN`
- `DEFAULT_WORKSHOP_ID`, default `1`

## Память

В `docker-compose.yml` заданы лимиты:

- backend: `mem_limit: 512m`, `GOMEMLIMIT=460MiB`
- PostgreSQL: `mem_limit: 512m`
- frontend nginx: `mem_limit: 256m`
