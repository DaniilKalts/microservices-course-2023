# Week 1 Homework (gRPC)

**Goal:** build 2 gRPC services (`user`, `chat`) with `.proto` contracts, generated stubs, and minimal server implementations (log request → return dummy response).

---

## Status

### Common (both services)
- [x] Create repo + Go module
- [x] Add `proto/` directory with service definitions
- [x] Add codegen command/script (protoc + Go plugins)
- [x] Generate Go stubs from `.proto`
- [x] Add `cmd/<service>/main.go` to start a gRPC server
- [x] Register service implementation
- [x] Basic logging
- [x] README: generate / run / quick call examples

### Service: `user` (User API)
*Directory: `cmd/user`*

**Features**
- User CRUD over gRPC

**Proto**
- [x] `Role` enum: `user`, `admin`
- [x] `Create(name, email, password, password_confirm, role) -> (id)`
- [x] `Get(id) -> (id, name, email, role, created_at, updated_at)`
- [x] `Update(id, name?, email?) -> Empty` (partial update)
- [x] `Delete(id) -> Empty`

**Implementation**
- [x] Implement handlers: log request → return dummy response / `Empty`

### Service: `chat` (Chat API)
*Directory: `cmd/chat`*

**Features**
- Chat management + sending messages over gRPC

**Proto**
- [x] `Create(usernames[]) -> (id)`
- [x] `Delete(id) -> Empty`
- [x] `SendMessage(from, text, timestamp) -> Empty`

**Implementation**
- [x] Implement handlers: log request → return dummy response / `Empty`

---

## 🚀 How to Run

### Prerequisites
- Go 1.21+
- [Task](https://taskfile.dev/) (optional, for convenient generation)
- [grpcurl](https://github.com/fullstorydev/grpcurl) (for testing endpoints)

### 1. Generate Code
If you modify `.proto` files, regenerate the Go stubs:
```bash
task generate
```

### 2. Run Services
Run both services (Chat on port 50051, User on port 50052) in parallel with a single command:
```bash
task dev
```

### 3. Test with `grpcurl`

**Create a User:**
```bash
grpcurl -plaintext -d '{"name": "Gopher", "email": "gopher@go.dev", "role": "ADMIN"}' localhost:50052 user.v1.UserV1/Create
```

**Get a User:**
```bash
grpcurl -plaintext -d '{"id": "some-uuid"}' localhost:50052 user.v1.UserV1/Get
```

**Send a Message:**
```bash
grpcurl -plaintext -d '{"from": "Alice", "text": "Hello Gopher!"}' localhost:50051 chat.v1.ChatV1/SendMessage
```
