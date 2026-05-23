#!/usr/bin/env bash
set -euo pipefail

PROJECT_DIR="${1:-repair-crm}"

mkdir -p "$PROJECT_DIR"/{cmd/app,internal/{api,models,repository,service},pkg/{auth,hash,password},migrations,frontend/src/{api,components,pages},frontend/public,scripts}

touch "$PROJECT_DIR"/cmd/app/main.go
touch "$PROJECT_DIR"/internal/api/{router.go,middleware.go,response.go,auth_handler.go,ticket_handler.go}
touch "$PROJECT_DIR"/internal/models/models.go
touch "$PROJECT_DIR"/internal/repository/{db.go,errors.go,master_repository.go,ticket_repository.go}
touch "$PROJECT_DIR"/internal/service/{errors.go,auth_service.go,ticket_service.go}
touch "$PROJECT_DIR"/pkg/auth/jwt.go
touch "$PROJECT_DIR"/pkg/hash/hash.go
touch "$PROJECT_DIR"/pkg/password/password.go
touch "$PROJECT_DIR"/migrations/{000001_init.up.sql,000001_init.down.sql,000002_public_flow.up.sql}
touch "$PROJECT_DIR"/frontend/{package.json,tsconfig.json,tsconfig.node.json,vite.config.ts,tailwind.config.js,postcss.config.js,index.html,Dockerfile,nginx.conf}
touch "$PROJECT_DIR"/frontend/src/{App.tsx,index.css,main.tsx,vite-env.d.ts}
touch "$PROJECT_DIR"/frontend/src/api/client.ts
touch "$PROJECT_DIR"/frontend/src/components/{StatusBadge.tsx,ui.ts}
touch "$PROJECT_DIR"/frontend/src/pages/{Dashboard.tsx,PublicApplication.tsx,StatusPage.tsx}
touch "$PROJECT_DIR"/{Dockerfile,docker-compose.yml,Makefile,README.md,.gitignore,.env.example}

cd "$PROJECT_DIR"
git init
git config user.name "${GIT_USER_NAME:-Repair CRM Bot}"
git config user.email "${GIT_USER_EMAIL:-repair-crm@example.local}"
git add .
git commit -m "Initial empty Repair CRM layout"
