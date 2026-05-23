COMPOSE := docker compose
DB_SERVICE := postgres

.PHONY: run migrate migrate-down seed-demo build logs down clean backend-test frontend-build

run: migrate
	$(COMPOSE) up --build

migrate:
	$(COMPOSE) up -d $(DB_SERVICE)
	$(COMPOSE) exec -T $(DB_SERVICE) sh -c 'for i in $$(seq 1 30); do pg_isready -U "$${POSTGRES_USER}" -d "$${POSTGRES_DB}" && break; sleep 1; done; for f in /migrations/*.up.sql; do psql -v ON_ERROR_STOP=1 -U "$${POSTGRES_USER}" -d "$${POSTGRES_DB}" -f "$$f"; done'

migrate-down:
	$(COMPOSE) exec -T $(DB_SERVICE) sh -c 'psql -v ON_ERROR_STOP=1 -U "$${POSTGRES_USER}" -d "$${POSTGRES_DB}" -f /migrations/000001_init.down.sql'

seed-demo:
	$(COMPOSE) exec -T $(DB_SERVICE) psql -v ON_ERROR_STOP=1 -U repair -d repair_crm -c "INSERT INTO workshops (id, name) VALUES (1, 'Demo Repair') ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name; INSERT INTO masters (workshop_id, username, password_hash) VALUES (1, 'admin', '\$$2b\$$12\$$xxHRJWHOPwTGxH.a/0xKz.2fWSAaA7t7tzwyZHDsG0j5dr4UCB6rK') ON CONFLICT (username) DO UPDATE SET workshop_id = EXCLUDED.workshop_id, password_hash = EXCLUDED.password_hash;"

build:
	$(COMPOSE) build

logs:
	$(COMPOSE) logs -f --tail=100

down:
	$(COMPOSE) down

clean:
	$(COMPOSE) down -v --remove-orphans

backend-test:
	go test ./...

frontend-build:
	cd frontend && npm ci && npm run build
