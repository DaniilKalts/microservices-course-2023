# Week 4 Homework (User Service)

**Goal:** 

- - -

### Status

**To-do:** 
- [] 
- []
- []

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
