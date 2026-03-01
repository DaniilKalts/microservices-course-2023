# Week 7 Homework

**Goal:** Add observability to the service with a complete stack: structured logger, logging interceptor/middleware, Prometheus, Loki, Promtail, Grafana, and Jaeger.

- - -

### Status

**Observability To-do:**
- [x] Add shared structured logger for the application (JSON output, levels, service metadata)
- [x] Add gRPC logging interceptor for unary requests
- [x] Expose `/metrics` and register Prometheus metrics (requests, latency, errors)
- [x] Add `prometheus` service with scrape configuration in `docker-compose.yaml`
- [x] Add `alertmanager` service and wire Prometheus alerts
- [x] Add `grafana` service with provisioned Prometheus datasource
- [x] Add `node-exporter` service for host CPU/memory/saturation metrics
- [ ] Add `loki` service for centralized log storage
- [ ] Add `promtail` service to collect and ship logs to Loki
- [ ] Add `jaeger` service and tracing export from HTTP/gRPC flows
- [ ] Add Grafana data sources for Loki and Jaeger
- [ ] Add dashboards and a simple observability smoke-check guide

- - -

### Project Structure

```text
cmd/
  main.go
internal/
  adapters/
    in/
      database/postgres/
      transport/http/
        middleware/
        swagger/
    out/
      transport/grpc/
        handlers/
        interceptor/
  app/
  clients/database/
  config/
    env/
  domain/
    auth/
    user/
  repository/
    mocks/
    user/
  service/
    auth/
    user/
deployments/
  alertmanager/
    alertmanager.yml
    entrypoint.sh
  grafana/
    provisioning/
      datasources/
        datasources.yml
  jaeger/
  loki/
  migrations/
    00001_create_users_table.sql
    Dockerfile
    migrate.sh
  promtail/
  prometheus/
    alerts.yml
    prometheus.yml
proto/
  auth/v1/
  user/v1/
pkg/
  env/
  jwt/
Taskfile.yaml
Dockerfile
docker-compose.yaml
third_party/
  proto/
  swagger-ui/

# generated after `task generate`
gen/
  grpc/
  openapi/
build/
  jwt/
  tls/
  tests/
```

- - -

### How to Run

### Prerequisites
- Go 1.25+
- Docker and Docker Compose
- Task (`go-task`)

### 1) Run with Docker

```bash
cp .env.example .env
docker compose up -d --build
```

### Service URLs (Docker)

- API Gateway: `http://localhost:8000`
- Swagger UI: `http://localhost:8000/swagger/`
- App metrics endpoint: `http://localhost:2112/metrics`
- Node Exporter metrics endpoint: `http://localhost:9100/metrics`
- Prometheus: `http://localhost:9090`
- Alertmanager: `http://localhost:9093`
- Grafana: `http://localhost:3000` (credentials from `.env`, default `admin/admin`)
- Loki (planned): `http://localhost:3100`
- Jaeger UI (planned): `http://localhost:16686`

### 2) Run Locally (Go)

```bash
cp .env.example .env
docker compose up -d postgres
./deployments/migrations/migrate.sh
task generate
task run-user
```

- - -

### Generate Build/Gen Artifacts

```bash
task generate      # generate grpc/gateway/openapi, swagger statik, tls certs, jwt keys
task tls:generate  # generate only TLS cert/key
task jwt:generate  # generate only JWT key pair
```

- - -

### Cleanup

```bash
task clean:gen    # remove generated code (gen/* + swagger statik.go)
task clean:build  # remove build/runtime artifacts (build/*)
task clean        # run all cleanup tasks
```
