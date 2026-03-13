# 🎓 Go Microservices Course (2023)

I took an 8-week microservices course and built everything from the ground up — week by week, feature by feature. It starts with a simple gRPC service and by the end turns into a **real, production-ready User Management API** with authentication, monitoring, and everything you'd expect from a serious backend. This repo is basically my progress diary — each folder is a snapshot of what I learned and built that week. If you're curious how microservice patterns and tools work in a real Go project, feel free to explore!

---

### ✨ Features

- 🔗 **gRPC & REST API**
    - Protocol Buffers contract-first API design
    - gRPC-Gateway for automatic HTTP/JSON transcoding
    - Swagger UI for interactive API documentation
    - Request validation via protoc-gen-validate
- 👤 **User Management (Admin Only)**
    - Create, Read, Update, Delete users
    - List all users
    - Role-based access control (User / Admin)
- 🙋 **Profile Management (Authenticated Users)**
    - View own profile
    - Update name, email, or password
    - Delete own account
- 🔐 **Authentication & Security**
    - JWT-based authentication (RS256 key pair)
    - Register, Login, Logout, Token Refresh flows
    - Access & Refresh token rotation
    - TLS encryption for gRPC transport
    - Password strength validation (timing-attack resistant)
- 🛡️ **Resilience & Reliability**
    - Rate limiting with per-IP token bucket (separate limits for auth endpoints)
    - Circuit breaker for gateway → gRPC calls (fail-fast on cascading failures)
    - Request timeout management at interceptor level
    - Health check endpoints (liveness & readiness probes)
    - Centralized gRPC error handling with user-friendly messages
- 📊 **Observability Stack**
    - Structured JSON logging (Zap) with Loki aggregation
    - Prometheus metrics collection with alert rules
    - Alertmanager with Telegram notifications
    - Distributed tracing via Jaeger + Tempo
    - Grafana dashboards (Golden Signals)
    - Promtail log shipping, Node Exporter host metrics
- 🏗️ **Clean Architecture & Testing**
    - Layered architecture: Handler → Service → Repository
    - Domain-driven design with dedicated domain models
    - Dependency injection container
    - Transaction manager for atomic DB operations
    - Unit tests with mocks (minimock) and table-driven patterns
    - Integration tests for repository layer
    - Load testing with ghz
- 🐳 **Infrastructure & DevOps**
    - Multi-stage Docker builds (non-root runtime)
    - Docker Compose orchestration (14 services)
    - Automated database migrations (goose)
    - Taskfile-based workflows (generate, test, migrate, TLS/JWT key management)

---

### 🛠 Tech Stack

| Category | Technologies |
|---|---|
| **Language** | Go |
| **Transport** | gRPC, gRPC-Gateway (REST), Protocol Buffers |
| **Database** | PostgreSQL 17, pgx driver, Squirrel query builder |
| **Authentication** | JWT (RS256), TLS certificates |
| **Observability** | Prometheus, Grafana, Loki, Promtail, Jaeger, Tempo, Alertmanager |
| **Testing** | minimock, gotestsum, ghz (load testing) |
| **Containerization** | Docker, Docker Compose |
| **Migrations** | goose |
| **Task Runner** | Taskfile |
| **Logging** | Zap (structured JSON) |

---

### 📅 Weekly Progress

| Week | Theme | Key Additions |
|---|---|---|
| **1** | gRPC Foundation | Proto definitions, User & Chat gRPC services, code generation |
| **2** | Database & Docker | PostgreSQL, Docker Compose, migrations (goose), config management, CI/CD |
| **3** | Clean Architecture | Handler → Service → Repository layers, DI container, transaction manager |
| **4** | Testing | Unit tests with mocks, table-driven tests, coverage reports |
| **5** | HTTP Gateway & Docs | gRPC-Gateway, request validation, Swagger UI, OpenAPI spec generation |
| **6** | Security | TLS encryption, JWT auth (register/login/logout/refresh), auth interceptor |
| **7** | Observability | Prometheus, Grafana, Loki, Promtail, Jaeger, Alertmanager, structured logging |
| **8** | Resilience | Rate limiter, circuit breaker, timeouts, health checks, error handling |

---

### 📂 Project Structure

> The repository is organized by week. Each directory is a self-contained Go project
> representing the cumulative state of the system at that point in the course.
> Below, weeks 1–7 show only key highlights; **week 8** (the final state) is shown in full.

```
.
├── 1-week/                        # gRPC foundation
│   ├── cmd/{chat,user}/           #   Two separate service entry points
│   ├── proto/{chat,user}/v1/      #   Proto contract definitions
│   ├── gen/go/{chat,user}/v1/     #   Generated gRPC stubs
│   └── Taskfile.yaml
│
├── 2-week/                        # Database & Docker
│   ├── internal/config/env/       #   Environment-based configuration
│   ├── migrations/                #   SQL migration files (goose)
│   ├── docker-compose.yaml        #   Postgres + app services
│   ├── Dockerfile                 #   Multi-stage build
│   └── local.env                  #   Environment template
│
├── 3-week/                        # Clean architecture
│   ├── internal/
│   │   ├── api/grpc/{chat,user}/  #   gRPC handlers (transport layer)
│   │   ├── app/                   #   Bootstrap & DI container
│   │   ├── clients/database/      #   DB client, transaction manager
│   │   ├── converter/             #   Proto ↔ domain mappers
│   │   ├── models/                #   Domain models
│   │   ├── repository/user/       #   Data access layer
│   │   └── service/user/          #   Business logic layer
│   └── ...
│
├── 4-week/                        # Testing & refactoring
│   ├── internal/
│   │   ├── adapters/{in,out}/     #   Ports & adapters pattern
│   │   ├── domain/user/           #   Domain entities & value objects
│   │   └── service/user/*_test.go #   Unit tests with mocks
│   └── ...
│
├── 5-week/                        # HTTP gateway & validation
│   ├── internal/adapters/
│   │   ├── in/transport/http/     #   Swagger UI handler
│   │   └── out/transport/grpc/
│   │       ├── interceptor/       #   Request validation
│   │       └── user/list.go       #   New List endpoint
│   ├── gen/openapi/               #   Generated OpenAPI specs
│   ├── third_party/swagger-ui/    #   Swagger UI static assets
│   └── ...
│
├── 6-week/                        # Security (TLS + JWT)
│   ├── build/{jwt,tls}/           #   Key pairs & certificates
│   ├── internal/
│   │   ├── adapters/out/transport/grpc/
│   │   │   ├── handlers/auth/     #   Auth handler (login, register, etc.)
│   │   │   └── interceptor/auth.go#   JWT auth interceptor
│   │   ├── domain/auth/           #   Auth domain (credentials, tokens)
│   │   └── service/auth/          #   Auth business logic
│   ├── pkg/jwt/                   #   JWT manager (RS256)
│   └── ...
│
├── 7-week/                        # Observability
│   ├── deployments/observability/ #   Full monitoring stack configs
│   │   ├── alertmanager/          #     Alert routing & Telegram
│   │   ├── grafana/               #     Dashboards & datasources
│   │   ├── jaeger/                #     Distributed tracing
│   │   ├── loki/                  #     Log aggregation
│   │   ├── prometheus/            #     Metrics & alert rules
│   │   ├── promtail/              #     Log shipping
│   │   └── tempo/                 #     Trace backend
│   ├── internal/adapters/transport/
│   │   ├── grpc/interceptor/
│   │   │   ├── logging.go         #   Request logging interceptor
│   │   │   ├── metrics.go         #   Prometheus metrics interceptor
│   │   │   └── tracing.go         #   Span creation interceptor
│   │   └── http/metrics/          #   Metrics HTTP server
│   ├── pkg/{logger,tracing}/      #   Shared observability packages
│   └── ...
│
├── 8-week/                        # Resilience & final state (full tree below)
│   ├── api/
│   │   ├── gen/
│   │   │   ├── go/
│   │   │   │   ├── auth/v1/                    # Generated auth gRPC stubs
│   │   │   │   │   ├── auth.pb.go
│   │   │   │   │   ├── auth.pb.gw.go
│   │   │   │   │   ├── auth.pb.validate.go
│   │   │   │   │   └── auth_grpc.pb.go
│   │   │   │   └── user/v1/                    # Generated user gRPC stubs
│   │   │   │       ├── profile.pb.go
│   │   │   │       ├── profile.pb.gw.go
│   │   │   │       ├── profile.pb.validate.go
│   │   │   │       ├── profile_grpc.pb.go
│   │   │   │       ├── user.pb.go
│   │   │   │       ├── user.pb.gw.go
│   │   │   │       ├── user.pb.validate.go
│   │   │   │       └── user_grpc.pb.go
│   │   │   └── openapi/                        # Generated OpenAPI specs
│   │   │       ├── auth/v1/auth.swagger.json
│   │   │       ├── user/v1/
│   │   │       │   ├── profile.swagger.json
│   │   │       │   └── user.swagger.json
│   │   │       └── gateway.swagger.json
│   │   ├── proto/                               # Proto contract definitions
│   │   │   ├── auth/v1/auth.proto
│   │   │   └── user/v1/
│   │   │       ├── profile.proto
│   │   │       └── user.proto
│   │   └── third_party/                         # Third-party proto deps
│   │       ├── google/api/
│   │       ├── protoc-gen-openapiv2/options/
│   │       └── validate/
│   ├── build/
│   │   ├── jwt/                                 # RS256 key pair
│   │   │   ├── rs256_private.pem
│   │   │   └── rs256_public.pem
│   │   └── tls/                                 # TLS certificates
│   │       ├── server.crt
│   │       └── server.key
│   ├── cmd/
│   │   └── main.go                              # Application entry point
│   ├── deploy/
│   │   ├── migrations/                          # Database migrations
│   │   │   ├── 00001_create_users_table.sql
│   │   │   ├── Dockerfile
│   │   │   └── migrate.sh
│   │   ├── observability/                       # Monitoring stack configs
│   │   │   ├── alertmanager/
│   │   │   ├── grafana/
│   │   │   │   ├── dashboards/golden-signals.json
│   │   │   │   └── provisioning/
│   │   │   ├── jaeger/jaeger.yml
│   │   │   ├── loki/loki.yml
│   │   │   ├── prometheus/
│   │   │   │   ├── alerts.yml
│   │   │   │   └── prometheus.yml
│   │   │   ├── promtail/promtail.yml
│   │   │   └── tempo/tempo.yml
│   │   └── scripts/                             # Server setup scripts
│   │       ├── install-docker-ubuntu.sh
│   │       └── setup-server.sh
│   ├── internal/
│   │   ├── adapters/
│   │   │   ├── database/postgres/               # PostgreSQL adapter
│   │   │   │   ├── client.go
│   │   │   │   ├── db.go
│   │   │   │   └── errors.go
│   │   │   └── transport/
│   │   │       ├── grpc/
│   │   │       │   ├── handlers/
│   │   │       │   │   ├── auth/                # Auth gRPC handler
│   │   │       │   │   │   ├── converter.go
│   │   │       │   │   │   └── handler.go
│   │   │       │   │   ├── profile/             # Profile gRPC handler
│   │   │       │   │   │   ├── converter.go
│   │   │       │   │   │   └── handler.go
│   │   │       │   │   └── user/                # User admin gRPC handler
│   │   │       │   │       ├── converter.go
│   │   │       │   │       └── handler.go
│   │   │       │   ├── interceptor/             # gRPC interceptors
│   │   │       │   │   ├── auth/
│   │   │       │   │   │   ├── interceptor.go
│   │   │       │   │   │   ├── policy.go
│   │   │       │   │   │   ├── policy_default.go
│   │   │       │   │   │   └── policy_default_test.go
│   │   │       │   │   ├── errors.go
│   │   │       │   │   ├── logging.go
│   │   │       │   │   ├── metrics.go
│   │   │       │   │   ├── ratelimit.go
│   │   │       │   │   ├── timeout.go
│   │   │       │   │   ├── tracing.go
│   │   │       │   │   └── validation.go
│   │   │       │   └── server.go
│   │   │       └── http/
│   │   │           ├── diagnostic/              # Health & metrics server
│   │   │           │   ├── checkers.go
│   │   │           │   └── server.go
│   │   │           ├── gateway/                 # gRPC-Gateway (REST proxy)
│   │   │           │   ├── interceptor/
│   │   │           │   │   ├── circuitbreaker.go
│   │   │           │   │   └── tracing.go
│   │   │           │   ├── middleware/tracing.go
│   │   │           │   ├── proxy.go
│   │   │           │   └── swagger.go
│   │   │           └── swagger/handler.go
│   │   ├── app/                                 # Bootstrap & DI
│   │   │   ├── app.go
│   │   │   └── container.go
│   │   ├── clients/database/                    # DB abstraction
│   │   │   ├── prettier/prettier.go
│   │   │   ├── transaction/transaction.go
│   │   │   └── database.go
│   │   ├── config/                              # Configuration
│   │   │   ├── config.go
│   │   │   └── load.go
│   │   ├── domain/user/                         # Domain layer
│   │   │   ├── credentials.go
│   │   │   ├── errors.go
│   │   │   ├── role.go
│   │   │   └── user.go
│   │   ├── repository/user/                     # Data access layer
│   │   │   ├── converter.go
│   │   │   ├── logging.go
│   │   │   ├── mock.go
│   │   │   ├── model.go
│   │   │   ├── repository.go
│   │   │   ├── repository_integration_test.go
│   │   │   └── tracing.go
│   │   └── service/                             # Business logic
│   │       ├── auth/
│   │       │   ├── errors.go
│   │       │   ├── logging.go
│   │       │   ├── service.go
│   │       │   ├── service_test.go
│   │       │   ├── tracing.go
│   │       │   └── types.go
│   │       ├── user/
│   │       │   ├── logging.go
│   │       │   ├── mock.go
│   │       │   ├── service.go
│   │       │   ├── service_test.go
│   │       │   └── tracing.go
│   │       └── services.go
│   ├── pkg/                                     # Shared packages
│   │   ├── env/parse.go
│   │   ├── jwt/
│   │   │   ├── claims.go
│   │   │   ├── errors.go
│   │   │   ├── keys.go
│   │   │   ├── manager.go
│   │   │   ├── mock.go
│   │   │   └── token.go
│   │   ├── logger/logger.go
│   │   ├── protoutil/wrappers.go
│   │   └── tracing/tracer.go
│   ├── tasks/                                   # Taskfile includes
│   │   ├── clean.yaml
│   │   ├── generate.yaml
│   │   ├── jwt.yaml
│   │   ├── migration.yaml
│   │   ├── test.yaml
│   │   └── tls.yaml
│   ├── web/swagger-ui/                          # Swagger UI assets
│   ├── docker-compose.yaml
│   ├── Dockerfile
│   ├── Taskfile.yaml
│   └── go.mod
└── ...
```

---

### 📡 API Endpoints

#### Auth Service (`auth.v1.AuthV1`)

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/v1/auth/register` | Create a new account |
| `POST` | `/api/v1/auth/login` | Authenticate and receive tokens |
| `POST` | `/api/v1/auth/logout` | Invalidate refresh token |
| `POST` | `/api/v1/auth/refresh` | Rotate access & refresh tokens |

#### Profile Service (`user.v1.ProfileV1`)

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/api/v1/users/me` | Get authenticated user's profile |
| `PATCH` | `/api/v1/users/me` | Update own name, email, or password |
| `DELETE` | `/api/v1/users/me` | Delete own account |

#### User Admin Service (`user.v1.UserV1`)

| Method | Endpoint | Description |
|---|---|---|
| `POST` | `/api/v1/users` | Create a managed user _(admin only)_ |
| `GET` | `/api/v1/users` | List all users |
| `GET` | `/api/v1/users/{id}` | Get user by ID |
| `PATCH` | `/api/v1/users/{id}` | Update a user _(admin only)_ |
| `DELETE` | `/api/v1/users/{id}` | Delete a user _(admin only)_ |

> **Access policy:** Register, Login, Refresh, List users, and Get user are **public** (no token required).
> Profile endpoints and Logout require **authentication** (any role). Create, Update, and Delete user are **admin only**.

#### Diagnostic Endpoints (port `2112`)

| Method | Endpoint | Description |
|---|---|---|
| `GET` | `/healthz/liveness` | Always returns `200 OK` — is the process alive? |
| `GET` | `/healthz/readiness` | Checks dependencies (e.g. database ping); returns `503` if any fail |
| `GET` | `/metrics` | Prometheus metrics scrape endpoint |

---

### 🌐 Web UIs

| Service | URL | Description |
|---|---|---|
| **Swagger UI** | `http://localhost:8000/api/swagger/` | Interactive REST API documentation & testing |
| **Grafana** | `http://localhost:3000` | Dashboards for metrics, logs, and traces (Golden Signals) |
| **Prometheus** | `http://localhost:9090` | Raw metrics queries and alert rule status |
| **Alertmanager** | `http://localhost:9093` | Active alerts, silences, and routing config |
| **Jaeger** | `http://localhost:16686` | Distributed trace search and visualization |
| **Dozzle** | `http://localhost:8080` | Real-time Docker container log viewer |

---

### 🖼️ Screenshots

#### Swagger UI — Interactive API Documentation

![Swagger UI](screenshots/microservices_course_swagger.png)

#### Jaeger — Trace Search (GET /api/v1/users)

![Jaeger Trace Search](screenshots/microservices_course_jaeger_traces.png)

#### Jaeger — Trace Detail with Span Waterfall & Error Logs

![Jaeger Trace Detail](screenshots/microservices_course_jaeger_trace.png)

#### Grafana — Golden Signals Dashboard (CPU, Memory, Latency, Errors, Request Rate)

![Grafana Golden Signals — Part 1](screenshots/microservices_course_grafana_1.png)

#### Grafana — Memory Breakdown, API Logs (Loki) & Recent Traces (Tempo)

![Grafana Golden Signals — Part 2](screenshots/microservices_course_grafana_2.png)

#### Dozzle — Real-Time Docker Container Logs

![Dozzle](screenshots/microservices_course_dozzle.png)

---

### 🏗 Setup & Installation

#### Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- [Docker & Docker Compose](https://docs.docker.com/engine/install/)
- [Task](https://taskfile.dev/installation/) (task runner)

#### 1. Clone the Repository

```bash
git clone https://github.com/DaniilKalts/microservices-course-2023.git
cd microservices-course-2023/8-week
```

#### 2. Configure Environment Variables

Create a `.env` file in the `8-week/` directory. Below is a complete reference with default values:

```bash
# ─── Application ──────────────────────────────────────────────
APP_SHUTDOWN_TIMEOUT=5s              # Graceful shutdown window

# ─── PostgreSQL ───────────────────────────────────────────────
POSTGRES_DB=postgres                 # Database name
POSTGRES_USER=postgres               # Database user
POSTGRES_PASSWORD=postgres           # Database password
POSTGRES_HOST=postgres               # Hostname (use "postgres" for Docker, "localhost" for local)
POSTGRES_PORT=5432                   # Database port
POSTGRES_SSLMODE=disable            # SSL mode
POSTGRES_QUERY_TIMEOUT=1s           # Per-query timeout

# ─── gRPC Server ─────────────────────────────────────────────
GRPC_HOST=0.0.0.0                   # gRPC listen address
GRPC_PORT=50051                     # gRPC port
GRPC_REQUEST_TIMEOUT=5s             # Max request duration
GRPC_RATE_LIMIT_RPS=100             # Default rate limit (requests/sec)
GRPC_RATE_LIMIT_BURST=20            # Default burst size
GRPC_RATE_LIMIT_AUTH_RPS=5          # Auth endpoint rate limit (stricter)
GRPC_RATE_LIMIT_AUTH_BURST=10       # Auth endpoint burst size

# ─── HTTP Gateway ────────────────────────────────────────────
GATEWAY_HOST=0.0.0.0                # REST gateway listen address
GATEWAY_PORT=8000                   # REST gateway port
GATEWAY_CB_MAX_REQUESTS=3           # Circuit breaker: max requests in half-open
GATEWAY_CB_OPEN_TIMEOUT=30s         # Circuit breaker: wait before retry
GATEWAY_CB_FAILURE_THRESHOLD=5      # Circuit breaker: failures to trip

# ─── Diagnostic Server ───────────────────────────────────────
DIAGNOSTIC_HOST=0.0.0.0             # Metrics/health listen address
DIAGNOSTIC_PORT=2112                # Prometheus metrics port

# ─── JWT (RS256) ─────────────────────────────────────────────
JWT_ISS=auth-service                # Token issuer
JWT_SUB=access-token                # Token subject
JWT_AUD=chatting-clients            # Token audience
JWT_PRIVATE_KEY_FILE=build/jwt/rs256_private.pem
JWT_PUBLIC_KEY_FILE=build/jwt/rs256_public.pem
JWT_ACCESS_EXP=15m                  # Access token TTL
JWT_REFRESH_EXP=168h                # Refresh token TTL (7 days)
JWT_NBF=0s                          # Not-before offset
JWT_IAT=0s                          # Issued-at offset

# ─── TLS ─────────────────────────────────────────────────────
TLS_ENABLED=false                   # Enable TLS for gRPC
TLS_CERT_FILE=build/tls/server.crt
TLS_KEY_FILE=build/tls/server.key

# ─── Logging (Zap) ───────────────────────────────────────────
ZAP_LEVEL=info                      # Log level (debug, info, warn, error)
ZAP_ENCODING=json                   # Output format (json, console)
ZAP_OUTPUT_PATHS=stdout             # Log output destination
ZAP_ERROR_OUTPUT_PATHS=stderr       # Error log destination

# ─── Tracing (Jaeger) ────────────────────────────────────────
TRACING_ENABLED=true                # Enable distributed tracing
TRACING_SERVICE_NAME=api            # Service name in traces
TRACING_JAEGER_AGENT_HOST=jaeger    # Jaeger agent host ("jaeger" for Docker)
TRACING_JAEGER_AGENT_PORT=6831      # Jaeger agent UDP port
TRACING_SAMPLER_TYPE=const          # Sampling strategy
TRACING_SAMPLER_PARAM=1             # Sample all requests (1 = 100%)

# ─── Alertmanager ─────────────────────────────────────────────
ALERTMANAGER_TELEGRAM_BOT_TOKEN=    # Telegram bot token for alerts
ALERTMANAGER_TELEGRAM_CHAT_ID=      # Telegram chat ID for alerts

# ─── Grafana ─────────────────────────────────────────────────
GRAFANA_ADMIN_USER=admin            # Grafana admin username
GRAFANA_ADMIN_PASSWORD=admin        # Grafana admin password
```

#### 3. Generate JWT & TLS Keys

```bash
# Generate RS256 key pair for JWT signing
task jwt:generate

# Generate self-signed TLS certificate (optional, if TLS_ENABLED=true)
task tls:generate
```

#### 4. Start the Project

```bash
docker compose up --build
```

This starts **14 services**: PostgreSQL, API server, migrator, Prometheus, Grafana, Loki, Promtail, Jaeger, Tempo, Alertmanager, Node Exporter, and Dozzle.

#### 5. Verify the Setup

Once all containers are running, the following endpoints become available:

| Service | URL | Description |
|---|---|---|
| **REST API** | `http://localhost:8000` | HTTP Gateway (JSON/REST) |
| **Swagger UI** | `http://localhost:8000/api/swagger/` | Interactive API docs |
| **gRPC** | `localhost:50051` | gRPC server (use grpcurl or ghz) |
| **Grafana** | `http://localhost:3000` | Dashboards (admin/admin) |
| **Prometheus** | `http://localhost:9090` | Metrics & queries |
| **Alertmanager** | `http://localhost:9093` | Alert management |
| **Jaeger** | `http://localhost:16686` | Distributed traces |
| **Dozzle** | `http://localhost:8080` | Docker container logs |

---

### ⚙️ Available Task Commands

All commands should be run from the `8-week/` directory.

```bash
# ─── Code Generation ─────────────────────────────────────────
task generate:generate     # Generate protobuf, gRPC, gateway, and OpenAPI stubs

# ─── Database Migrations ─────────────────────────────────────
task db:status             # Check current migration status
task db:new NAME=<name>    # Create a new SQL migration file
task db:up                 # Apply all pending migrations
task db:down               # Rollback the last migration

# ─── Testing ─────────────────────────────────────────────────
task test:test                # Run all unit tests
task test:test-coverage-html  # Run tests with HTML coverage report
task test:list-users          # Load test: List users endpoint (ghz)
task test:delete-user-fail    # Load test: Delete user — expected failure (ghz)

# ─── Security Keys ───────────────────────────────────────────
task jwt:generate          # Generate RS256 private/public key pair
task jwt:verify            # Verify JWT key pair integrity
task tls:generate          # Generate TLS certificate and key
task tls:verify            # Verify TLS certificate validity

# ─── Cleanup ─────────────────────────────────────────────────
task clean:gen             # Remove generated protobuf artifacts
task clean:build           # Remove build/runtime artifacts
task clean:all             # Remove all generated and build artifacts
```

---

### 🏃 Running Locally (Without Docker)

#### Prerequisites

Make sure the following are installed on your machine:

- [Go 1.25+](https://go.dev/dl/)
- [Task](https://taskfile.dev/installation/) — task runner
- [PostgreSQL](https://www.postgresql.org/download/) — running on `localhost:5432`
- [goose](https://github.com/pressly/goose) — migration tool (`go install github.com/pressly/goose/v3/cmd/goose@latest`)

> You can also run `task install-deps` to install goose and all other Go-based tools at once.

#### 1. Create the PostgreSQL Database

```bash
# Connect to PostgreSQL and create the database
psql -U postgres -c "CREATE DATABASE postgres;"
```

> If you use a different database name, user, or password — update the `POSTGRES_*` variables in `.env` accordingly.

#### 2. Create and Configure `.env`

Create a `.env` file in the `8-week/` directory. The key differences from the Docker setup:

```bash
# ─── PostgreSQL (point to localhost instead of Docker container) ───
POSTGRES_HOST=localhost

# ─── gRPC Server ──────────────────────────────────────────────────
GRPC_HOST=0.0.0.0
GRPC_PORT=50051
GRPC_REQUEST_TIMEOUT=5s

# ─── HTTP Gateway ─────────────────────────────────────────────────
GATEWAY_HOST=0.0.0.0
GATEWAY_PORT=8000

# ─── Diagnostic Server ───────────────────────────────────────────
DIAGNOSTIC_HOST=0.0.0.0
DIAGNOSTIC_PORT=2112

# ─── Disable tracing (Jaeger is not running locally) ─────────────
TRACING_ENABLED=false

# ─── Logging ─────────────────────────────────────────────────────
ZAP_LEVEL=info
ZAP_ENCODING=console
```

All other variables (`POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `POSTGRES_PORT`, `POSTGRES_SSLMODE`, `POSTGRES_QUERY_TIMEOUT`, `JWT_*`, `TLS_*`, rate limit, circuit breaker, etc.) can be copied from the [full `.env` reference](#2-configure-environment-variables) above. The defaults will work for most local setups.

> **Important:** `GRPC_HOST`, `GRPC_PORT`, `GRPC_REQUEST_TIMEOUT`, `GATEWAY_HOST`, `GATEWAY_PORT`, `DIAGNOSTIC_HOST`, `DIAGNOSTIC_PORT`, `JWT_ISS`, `JWT_SUB`, `JWT_AUD`, `JWT_ACCESS_EXP`, `JWT_REFRESH_EXP`, `JWT_NBF`, `JWT_IAT`, and all `POSTGRES_*` fields are **required** — the application will fail to start if they are missing.

#### 3. Generate JWT Keys

The application validates that the RS256 key pair files exist on startup. Generate them with:

```bash
task jwt:generate
```

This creates `build/jwt/rs256_private.pem` and `build/jwt/rs256_public.pem`.

#### 4. Run Database Migrations

The Taskfile reads `POSTGRES_*` variables from your `.env` to build the migration DSN automatically:

```bash
task db:up
```

Verify the migration was applied:

```bash
task db:status
```

#### 5. Start the Application

```bash
go run cmd/main.go -config-path=.env
```

The app will start three servers:
- **gRPC** on `:50051`
- **REST Gateway** on `:8000` (proxies to gRPC)
- **Diagnostic** on `:2112` (Prometheus metrics + health checks)

#### 6. Verify It Works

```bash
# Health check
curl http://localhost:2112/healthz/liveness

# Register a user via REST
curl -X POST http://localhost:8000/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"name": "John", "email": "john@example.com", "password": "Secret123"}'

# List users via REST
curl http://localhost:8000/api/v1/users
```

> **Optional:** If you want distributed tracing locally, start Jaeger separately (`docker run -p 16686:16686 -p 6831:6831/udp jaegertracing/jaeger:2.15.1`) and set `TRACING_ENABLED=true`, `TRACING_JAEGER_AGENT_HOST=localhost` in your `.env`.
