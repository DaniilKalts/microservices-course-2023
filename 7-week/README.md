# Week 7 Homework

**Goal:** Add observability to the service with a complete stack: structured logger, logging interceptor/middleware, Prometheus, Loki, Promtail, Grafana, and Jaeger.

---

### To-do

- [x] Add shared structured logger for the application (JSON output, levels, service metadata)
- [x] Add gRPC logging interceptor for unary requests
- [x] Expose `/metrics` and register Prometheus metrics (requests, latency, errors)
- [x] Add `prometheus` service with scrape configuration in `docker-compose.yaml`
- [x] Add `alertmanager` service and wire Prometheus alerts
- [x] Add `grafana` service with provisioned Prometheus datasource
- [x] Add `node-exporter` service for host CPU/memory/saturation metrics
- [x] Add `loki` service for centralized log storage
- [x] Add `promtail` service to collect and ship logs to Loki
- [x] Add `jaeger` service and tracing export from HTTP/gRPC flows
- [x] Add Grafana datasource for Loki
- [x] Add Grafana datasource for Jaeger

---

### Prerequisites

- Go 1.25+
- Docker and Docker Compose
- Task (`go-task`)

### Run

```bash
cp .env.example .env
docker compose up -d --build
```

---

### Web UI URLs

| Service             | URL                                |
| ------------------- | ---------------------------------- |
| API Gateway         | http://localhost:8000              |
| Swagger UI          | http://localhost:8000/swagger/     |
| Dozzle (Docker logs)| http://localhost:8080              |
| Grafana             | http://localhost:3000              |
| Prometheus          | http://localhost:9090              |
| Alertmanager        | http://localhost:9093              |
| Jaeger UI           | http://localhost:16686             |

Grafana credentials come from `.env` (default `admin` / `admin`).

### Internal Endpoints

| Service              | URL                                |
| -------------------- | ---------------------------------- |
| gRPC server          | localhost:50051                    |
| App metrics          | http://localhost:2112/metrics      |
| Node Exporter        | http://localhost:9100/metrics      |
| Loki API             | http://localhost:3100              |
| Tempo API            | http://localhost:3200              |
| Promtail             | http://localhost:9080              |

---

### Code Generation

```bash
task generate      # grpc/gateway/openapi, swagger statik, tls certs, jwt keys
task tls:generate  # TLS cert/key only
task jwt:generate  # JWT key pair only
```

### Cleanup

```bash
task clean:gen    # remove generated code (gen/* + swagger statik.go)
task clean:build  # remove build/runtime artifacts (build/*)
task clean        # run all cleanup tasks
```

---

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
  migrations/
    00001_create_users_table.sql
    Dockerfile
    migrate.sh
  observability/
    alertmanager/
    grafana/
    jaeger/
    loki/
    promtail/
    prometheus/
    tempo/
proto/
  auth/v1/
  user/v1/
pkg/
  env/
  jwt/
third_party/
  proto/
  swagger-ui/
Taskfile.yaml
Dockerfile
docker-compose.yaml

# generated after `task generate`
gen/
  grpc/
  openapi/
build/
  jwt/
  tls/
  tests/
```
