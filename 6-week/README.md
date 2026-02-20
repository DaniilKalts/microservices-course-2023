# Week 6 Homework (User Service)

**Goal:** Extend the **User** service with production-ready security features while preserving existing transport-level integrations (gRPC Gateway, validation, and Swagger UI).

- - -

### Status

**To-do:**
- [x] Add TLS encryption for gRPC server
- [x] Add JWT authentication with login, register, logout

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
  Dockerfile
  migrate.sh
proto/
  user/v1/user.proto
gen/
  grpc/user/v1/
  openapi/user/v1/
build/
  tests/
  tls/
Taskfile.yaml
Dockerfile
docker-compose.yaml
install-docker-ubuntu.sh
third_party/
  proto/
  swagger-ui/
```

- - -

### Implemented Modules

- **gRPC Gateway** exposes HTTP endpoints mapped to gRPC methods.
- **gRPC Validation** enforces request constraints before business logic execution.
- **Swagger UI** provides interactive API documentation for gateway endpoints.

### Planned Security Modules

- **TLS for gRPC** to encrypt service-to-service communication and harden transport security.
- **JWT Authentication** to secure access with register, login, and logout flows.

- - -

### How to Run

### Prerequisites
- Go 1.25+
- Docker and Docker Compose

### 1. Run Infrastructure and Services (Docker)

```bash
docker compose up -d --build
```

Swagger UI: `http://localhost:8000/swagger/`

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

### Cleanup Tasks

```bash
task clean:gen    # remove generated source (gen/grpc, gen/openapi, statik.go)
task clean:build  # remove build artifacts (build/tls, build/tests)
task clean        # run all cleanup tasks
```
