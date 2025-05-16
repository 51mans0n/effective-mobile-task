![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Postgres](https://img.shields.io/badge/postgres-%23316192.svg?style=for-the-badge&logo=postgresql&logoColor=white)
![Swagger](https://img.shields.io/badge/-Swagger-%23Clojure?style=for-the-badge&logo=swagger&logoColor=white)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=for-the-badge&logo=docker&logoColor=white)
![GitHub Actions](https://img.shields.io/badge/github%20actions-%232671E5.svg?style=for-the-badge&logo=githubactions&logoColor=white)

# Effective Mobile – Test Task
REST API for people enrichment (Go 1.23 + PostgreSQL + Docker)


> Creates a **person**, enriches it in parallel via  
> • [Agify.io](https://api.agify.io) – age  
> • [Genderize.io](https://api.genderize.io) – gender  
> • [Nationalize.io](https://api.nationalize.io) – probable country  
> …persists everything to PostgreSQL and exposes full CRUD.

---

## ✨ Features
|  Feature                                  | Description                                                                                                  |
|-------------------------------------------|--------------------------------------------------------------------------------------------------------------|
| 🗄 **Clean Architecture**                  | Clear separation `cmd/` (entry) → `internal/{client,service,repository,http}`. Easy support and testing.    |
| 🔗 **Parallel Enrichment**                | Simultaneous calls **agify**, **genderize**, **nationalize** through `errgroup`; 24 h in-memory cache.       |
| 🐘 **PostgreSQL + Migrations (embed)**     | SQL files are embedded into the binary (`//go:embed`), auto-migration on startup (`AUTO_MIGRATE=true`).                      |
| 📑 **Swagger UI**                          | Full OpenAPI 2.0; available at `/swagger/index.html`.                                                       |
| 🚦 **Graceful Shutdown**                   | Catches SIGINT/SIGTERM, correctly closes connections in ≤ 5 sec.                                             |
| 📋 **Structured Logging (zap)**           | Uniform JSON log: method, path, status, latency ms.                                                   |
| 🛡 **golangci-lint clean**                 | 0 issues; style / vet / errcheck / revive / staticcheck.                                                     |
| 🧪 **Unit-tests + sqlmock**                | The repository layer (Create / List / Update / Delete) and business logic are covered.                                    |
| 🐳 **Multi-stage Docker**                  | Builds a statically-linked binary, final image < 20 MB (distroless).                                  |
| ⚙️ **One-command run**                    | `docker compose up` — raises API + Postgres, everything is configured via `.env`.                            |

---

## 🛠 Stack
* **Go** 1.23
* **PostgreSQL** 16
* **chi** router • **sqlx** • **Squirrel** • **zap**
* **Swaggo** (`/docs`) 
* **golangci-lint**
* **Docker Compose v3.9**

---

## 🚀 Quick Start (Dev / Prod)

```bash
# 1. clone + prepare env
git clone https://github.com/51mans0n/effective-mobile-task.git
cd effective-mobile-task
cp .env.example .env      # edit if needed

# 2. build & run ↴  http://localhost:8080/swagger/index.html
docker compose up --build

# stop
docker compose down -v
```

### Local run without Docker
```bash
go run ./cmd/server      # needs running Postgres (see .env)
```

---

## 📂 Project Structure

```
.
├── cmd/
│   └── service/           → main.go (entry)
├── docs/                  → Swagger (generated)
├── internal/
│   ├── client/            → agify / genderize / nationalize + cache
│   ├── service/           → business logic
│   ├── repository/        → sqlx + Squirrel
│   ├── http/      
│   │   ├── handler/       → CRUD handlers
│   │   └── middleware/    → zap logger    
│   ├── migration/         → embedded .sql
│   ├── model/             → Person / Enriched
│   ├── config/            → .env loader
│   └── logger/            → zap factory  
├── .env.example           → example of environments
├── docker‑compose.yml     → docker-compose
├── Dockerfile             → Dockerfile
└── .golangci.yml          → lint
```

---

## 📜 API Overview

| Method   | Path           | Description                                                        |
| -------- | -------------- | ------------------------------------------------------------------ |
| `POST`   | `/people`      | Create person + enrichment                                         |
| `GET`    | `/people`      | List with filters `name`, `gender`, `country`, paging `page/limit` |
| `GET`    | `/people/{id}` | Get by UUID                                                        |
| `PUT`    | `/people/{id}` | Update F-I-O                                                       |
| `DELETE` | `/people/{id}` | Remove                                                             |

See Swagger UI for full schema & examples.

---

## 🔌 Environment Variables

| Var            | Default                                                              | Purpose                              |
| -------------- | -------------------------------------------------------------------- | ------------------------------------ |
| `APP_PORT`     | `:8080`                                                              | HTTP bind address                    |
| `LOG_LEVEL`    | `info`                                                               | debug / info / warn / error          |
| `AUTO_MIGRATE` | `true`                                                               | Run SQL migrations on startup        |
| `DB_DSN`       | `postgres://youruser:yourpass@postgres:5432/yourdb?sslmode=disable` | Postgres DSN                         |
| `CACHE_TTL`    | `24h`                                                                | In-memory cache TTL for external API |

---

## 🧪 Tests & Lint

```bash
go test ./...                 # unit tests
golangci-lint run ./...       # 0 issues
```

---

## 🙌 Author
- Maxim Skorokhod · Almaty
- GitHub https://github.com/51mans0n
- Telegram: Simanson