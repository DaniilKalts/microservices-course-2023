# Week 3 Homework (Clean Architecture, DI, Transaction Manager)

**Goal:** Refactor the **User** service to follow Clean Architecture principles (Handler-Service-Repository), implement a Dependency Injection (DI) container, and add a Transaction Manager.

---

## Status

### Service: `user`
*Directory: `cmd/user`, `internal`*

**Architecture & Patterns**
- [ ] **Clean Architecture:** Refactor to Handler-Service-Repository layers
    - Implement `Repository` layer (User CRUD)
    - Implement `Service` layer (Business logic)
    - Refactor `Handler` (gRPC) to use Service
- [ ] **DI Container:** Implement dependency injection container
- [ ] **Transaction Manager:** Add `TxManager` for atomic database operations

### Service: `chat`
*Directory: `cmd/chat`*
- No changes planned for this week.

---

## 🚀 How to Run

### Prerequisites
- Docker & Docker Compose
- [Goose](https://github.com/pressly/goose) (optional)
- [Task](https://taskfile.dev/) (optional)

### 1. Run Infrastructure & Services (Docker)
```bash
docker compose up -d --build
```

### 2. Run Locally (Go)

1. **Start Postgres:**
   ```bash
   docker compose up -d postgres
   ```

2. **Apply Migrations:**
   ```bash
   ./migrations/migrate.sh
   ```

3. **Run User Service:**
   ```bash
   go run cmd/user/main.go --config-path=local.env
   ```
