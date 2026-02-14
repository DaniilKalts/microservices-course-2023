# Week 4 Homework (User Service)

**Goal:** Finalize the **User** service implementation: keep clean architecture boundaries, and cover business logic with unit tests.

- - -

### Status

**To-do:** 
- [x] Organize the service using architecture layers (`adapters`, `app`, `domain`, `repository`, `service`)
- [x] Add unit tests for user service use cases (`create`, `get`, `update`, `delete`)
- [x] Add Task-based commands for test and coverage workflows

- - -

### Project Structure

```text
cmd/
  main.go
internal/
  adapters/
    in/database/postgres/
    out/transport/grpc/user/
  app/
  clients/database/
  config/
  domain/user/
  repository/
    repository.go
    user/
  service/
    service.go
    user/
      create.go
      create_test.go
      delete.go
      delete_test.go
      get.go
      get_test.go
      update.go
      update_test.go
migrations/
proto/
gen/
Taskfile.yaml
docker-compose.yaml
```

Notes:
- Unit tests currently live in `internal/service/user`.

- - -

### 🚀 How to Run

### Prerequisites
- Go 1.25+
- Docker and Docker Compose
- [Task](https://taskfile.dev/) (recommended)

### 1. Run Infrastructure & Services (Docker)

```bash
docker compose up -d --build
```

### 2. Run Locally (Go)

1. **Start PostgreSQL:**

```bash
docker compose up -d postgres
```

2. **Apply migrations:**

```bash
./migrations/migrate.sh
```

3. **Run service:**

```bash
go run ./cmd/main.go
```

- - -

### Testing With Task

Install tooling (protobuf plugins, minimock, gotestsum):

```bash
task install-deps
```

Run user service tests:

```bash
task test
```

Run user service tests with HTML coverage report:

```bash
task test-coverage-html
```

Coverage artifacts are generated under `.task/tests/` and are ignored by git.
