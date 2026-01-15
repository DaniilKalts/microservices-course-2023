# Week 2 Homework (Docker, DB, Config, CI/CD)

**Goal:** Containerize services, implement persistent storage with Postgres, manage configuration, handle migrations, and set up a CI/CD pipeline.

---

## Status

### Common / Infrastructure
- [x] **Configuration:** Load config from environment variables / file (`local.env`)
- [x] **Docker:** Add `Dockerfile` (multi-stage) for `user` and `chat` services
- [x] **Compose:** Add `docker-compose.yaml` for services, `postgres`, and `migrator`
- [x] **Migrations:** Add `migrator` service with `goose` and shell script

### Service: `user`
*Directory: `cmd/user`*

**Persistence**
- [x] **Migrations:** SQL schema for users table (`00001_create_users_table.sql`)
- [x] **Implementation:** Replace dummy handlers with `squirrel` SQL builder + `pgx` driver

### Service: `chat`
*Directory: `cmd/chat`*

### DevOps
- [x] **CI/CD:** GitHub Actions pipeline to build/push images and deploy to remote server via SSH

---

## 🚀 How to Run

### Prerequisites
- Docker & Docker Compose
- [Goose](https://github.com/pressly/goose) (optional, if running migrations manually)
- [Task](https://taskfile.dev/) (optional, for convenience)

### 1. Run Infrastructure & Services (Docker)
Start the database, apply migrations automatically, and run both services. The services are configured to use `local.env` via the `--config-path` flag.
```bash
docker compose up -d --build
```
> **Note:** The `migrator` container will wait for Postgres to be ready, apply any pending migrations, and then exit.

### 2. Run Locally (Go)
If you prefer running services outside Docker (e.g., for debugging):

1. **Start Postgres:**
   ```bash
   docker compose up -d postgres
   ```

2. **Apply Migrations:**
   ```bash
   # Using the provided script
   ./migrations/migrate.sh
   ```

3. **Run Services:**
   ```bash
   # User Service (Port 50052)
   go run cmd/user/main.go --config-path=local.env

   # Chat Service (Port 50051)
   go run cmd/chat/main.go --config-path=local.env
   ```
