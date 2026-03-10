# Week 8 Homework

**Goal:** Enhance service reliability and resilience by adding request timeouts, rate limiting, circuit breakers, health checks, and a custom error handler.

---

## What's Next

- [x] Add request timeout management.
- [ ] Implement rate limiter to prevent abuse.
- [ ] Add circuit breakers for external calls.
- [ ] Implement healthcheck patterns for service readiness and liveness.
- [ ] Create a custom error handler to provide user-friendly error messages.

---

## Prerequisites

- [Go 1.25+](https://go.dev/): To build and run the code.
- [Docker & Compose](https://www.docker.com/): To run the app and its dependencies.
- [Task](https://taskfile.dev/): To run predefined commands easily.

---

## Running the Application

Follow these steps to set up the environment, generate necessary code, and start the application:

1. **Set up the environment:**
```bash
cp .env.example .env
```

2. **Generate necessary code (gRPC, Gateway, OpenAPI, Swagger, etc.):**
```bash
task generate      # grpc/gateway/openapi, swagger statik, tls certs, jwt keys
task tls:generate  # TLS cert/key only
task jwt:generate  # JWT key pair only
```

*(Optional)* Before a fresh build, you can clean up previously generated artifacts:
```bash
task clean:gen    # remove generated code (gen/* + swagger statik.go)
task clean:build  # remove build/runtime artifacts (build/*)
task clean        # run all cleanup tasks
```

3. **Start the application and its dependencies:**
```bash
docker compose up -d --build
```

**Useful Docker Compose commands:**
- **Rebuild, remove old containers, and start new ones:** `docker compose up -d --build --force-recreate --remove-orphans`
- **View logs across all services:** `docker compose logs -f`
- **Stop and remove all services:** `docker compose down`
- **Stop, remove services and clean up volumes:** `docker compose down -v` (prevents memory and storage leaks)

---

## Web Interfaces

Explore the application observability and testing tools through your browser:

| Service | URL | Description |
|---|---|---|
| **API Gateway** | http://localhost:8000 | The main entry point for external HTTP traffic. |
| **Swagger UI** | http://localhost:8000/swagger/ | Interactive API documentation to easily test and understand the exposed HTTP endpoints. |
| **Dozzle** | http://localhost:8080 | A lightweight web interface to view real-time Docker container logs. |
| **Grafana** | http://localhost:3000 | The primary observability dashboard to visualize metrics, logs, and traces. Credentials come from `.env` (default `admin` / `admin`). |
| **Prometheus** | http://localhost:9090 | A monitoring system that scrapes and stores time-series metrics from our services. |
| **Alertmanager** | http://localhost:9093 | Handles alerts sent by Prometheus and routes them appropriately. |
| **Jaeger UI** | http://localhost:16686 | A distributed tracing UI to inspect how requests flow through different services. |

---

## Internal Services

These endpoints are mostly used internally by the application components and observability tools:

| Service | URL | Description |
|---|---|---|
| **gRPC server** | localhost:50051 | Internal gRPC server for high-performance communication between microservices. |
| **App metrics** | http://localhost:2112/metrics | The endpoint from which Prometheus scrapes the application's internal metrics. |
| **Node Exporter**| http://localhost:9100/metrics | Exposes hardware and OS metrics (like CPU and memory usage) to Prometheus. |
| **Loki API** | http://localhost:3100 | The backend endpoint where Promtail ships the logs for storage. |
| **Tempo API** | http://localhost:3200 | Backend endpoint for distributed tracing data storage. |
| **Promtail** | http://localhost:9080 | The agent that tails container logs and sends them to Loki. |

---



## Project Structure

```text
cmd/
  main.go
deployments/
  migrations/
  observability/
    alertmanager/
    grafana/
      dashboards/
      provisioning/
        dashboards/
        datasources/
    jaeger/
    loki/
    prometheus/
    promtail/
    tempo/
internal/
  adapters/
    database/
      postgres/
    transport/
      grpc/
        handlers/
          auth/
          profile/
          user/
        interceptor/
      http/
        gateway/
        metrics/
  app/
  clients/
    database/
      prettier/
      transaction/
  config/
  domain/
    user/
  repository/
    user/
  service/
    auth/
    user/
pkg/
  env/
  jwt/
  logger/
  protoutil/
  tracing/
proto/
  auth/v1/
  user/v1/
scripts/
tasks/
third_party/
  proto/
    google/api/
    protoc-gen-openapiv2/options/
    validate/
  swagger-ui/
Taskfile.yaml
Dockerfile
docker-compose.yaml
go.mod
go.sum

# generated after `task generate`
gen/
  grpc/
  openapi/
build/
  jwt/
  tls/
  tests/
```
