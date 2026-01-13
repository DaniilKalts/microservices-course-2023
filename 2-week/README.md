# Week 2 Homework (Docker, DB, Config, CI/CD)

**Goal:** Containerize services, implement persistent storage with Postgres, manage configuration, handle migrations, and set up a CI/CD pipeline.

---

## Status

### Common / Infrastructure
- [x] **Configuration:** Load config from environment variables / file
- [x] **Docker:** Add `Dockerfile` for `user` and `chat` services
- [x] **Compose:** Add `docker-compose.yaml` for services and Postgres
- [ ] **Migrations:** Add shell script for Goose migrations

### Service: `user`
*Directory: `cmd/user`*

**Persistence**
- [ ] **Migrations:** SQL schema for users
- [ ] **Implementation:** Replace dummy handlers with `squirrel` SQL builder + Postgres

### Service: `chat`
*Directory: `cmd/chat`*

**Persistence**
- [ ] **Migrations:** SQL schema for chats and messages
- [ ] **Implementation:** Replace dummy handlers with `squirrel` SQL builder + Postgres

### DevOps
- [x] **CI/CD:** Pipeline to build/push images and deploy to remote server

---

## 🚀 How to Run

### Prerequisites
- Go 1.21+
- Docker & Docker Compose
- [Goose](https://github.com/pressly/goose) (for migrations)
- [Task](https://taskfile.dev/) (optional)

### 1. Configuration
Ensure environment variables are set (or use `.env` file):
```bash
cp .env.example .env
```

### 2. Run Infrastructure (Postgres)
Start the database container:
```bash
docker-compose up -d db
```

### 3. Apply Migrations
Run the migration script to set up the database schema:
```bash
./migration.sh up
```

### 4. Run Services
**Local (Go):**
```bash
task dev
```

**Docker:**
```bash
docker-compose up --build
```
