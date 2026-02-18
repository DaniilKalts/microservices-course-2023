# Week 4 Homework (User Service)

**Goal:** Finalize the **User** service implementation with transport-level gRPC integrations for HTTP access, validation, and API documentation.

- - -

### Status

**To-do:**
- [x] Add gRPC Gateway for HTTP server integration
- [x] Add gRPC request validation
- [x] Add gRPC Swagger UI for API documentation

- - -

### Project Structure

```text
cmd/
  main.go
internal/
  adapters/
    in/database/postgres/
    in/transport/http/swagger/
    out/transport/grpc/
      interceptor/
      user/
  app/
  clients/database/
  config/
    env/
  domain/user/
  repository/
    mocks/
    repository.go
    user/
  service/
    service.go
    user/
      create.go
      delete.go
      get.go
      list.go
      update.go
migrations/
  00001_create_users_table.sql
  migrate.sh
proto/
  user/v1/user.proto
gen/
  go/user/v1/
  openapi/user/v1/
Taskfile.yaml
Dockerfile
docker-compose.yaml
third_party/
  proto/
  swagger-ui/
```

- - -

### Implemented Modules

- **gRPC Gateway** exposes HTTP endpoints mapped to gRPC methods.
- **gRPC Validation** enforces request constraints before business logic execution.
- **Swagger UI** provides interactive API documentation for the gateway endpoints.

- - -

### How to Run

### Prerequisites
- Go 1.25+
- Docker and Docker Compose

### 1. Run Infrastructure and Services (Docker)

```bash
docker compose up -d --build
```

### 2. Run Locally (Go)

1. Start PostgreSQL:

```bash
docker compose up -d postgres
```

2. Apply migrations:

```bash
./migrations/migrate.sh
```

3. Run the service:

```bash
go run ./cmd/main.go
```
