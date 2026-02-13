# Week 4 Homework (Unit Tests: User Service + API)

**Goal:** Add unit tests for the **User** service and **User** gRPC API layers.

- - -

### Status

#### Service: `user`
*Directory: `cmd/user`, `internal`*

**This week: Unit Tests**
- [x] **Service layer tests:** `internal/service/user`
  - Cover success + error paths (validation, repository errors)
  - Use mocks/fakes for repository and transaction manager
  - Prefer table-driven tests
- [ ] **API layer tests:** `internal/api/grpc/user`
  - Verify request mapping -> service calls -> response mapping
  - Cover service error propagation
  - Use a mocked `service.UserService`
- [ ] **Invoke unit tests:** run `go test ./...`
- [ ] **Invoke coverage:** run `go test ./... -coverprofile=coverage.out` and review `go tool cover -func=coverage.out`

- - -

### User Service Architecture (Per-Service Clean Architecture)

The `user` service now follows a service-first clean architecture layout with explicit layers:

```text
cmd/user/main.go
internal/
  adapters/
    in/database/postgres/
      client.go
      postgres.go
    out/transport/grpc/user/
      handler.go
      mapper.go
      create.go
      get.go
      update.go
      delete.go
  app/
    container.go
    app.go
  domain/user/
    entity.go
    patch.go
    role.go
  repository/user/
    repository.go
    model.go
    mapper.go
    create.go
    get.go
    update.go
    delete.go
  repository/repository.go
  service/service.go
  service/user/
    service.go
    create.go
    get.go
    update.go
    delete.go
```

Each use case and adapter method lives in its own file.

- - -

### How to Run

#### Prerequisites
- Docker & Docker Compose
- [Goose](https://github.com/pressly/goose) (optional)
- [Task](https://taskfile.dev/) (optional)

#### 1. Run Infrastructure & Service (Docker)
```bash
docker compose up -d --build
```

#### 2. Run Locally (Go)

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
   go run cmd/user/main.go --config-path=.env
   ```

#### 3. Run Unit Tests
```bash
go test ./...
```
