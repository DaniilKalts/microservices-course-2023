# Week 1 Homework Plan (gRPC)

**Goal:** build 2 gRPC services (`auth`, `chat-server`) with `.proto` contracts, generated stubs, and minimal server implementations (log request → return dummy response).

---

## Common (both services)
- [ ] Create repo + Go module
- [ ] Add `proto/` directory with service definitions
- [ ] Add codegen command/script (protoc + Go plugins)
- [ ] Generate Go stubs from `.proto`
- [ ] Add `cmd/<service>/main.go` to start a gRPC server
- [ ] Register service implementation
- [ ] Basic logging
- [ ] README: generate / run / quick call examples

---

## Service: `auth` (User API)

### Features
- User CRUD over gRPC

### Proto
- [ ] `Role` enum: `user`, `admin`
- [ ] `Create(name, email, password, password_confirm, role) -> (id)`
- [ ] `Get(id) -> (id, name, email, role, created_at, updated_at)`
- [ ] `Update(id, name?, email?) -> Empty` (partial update)
- [ ] `Delete(id) -> Empty`

### Implementation
- [ ] Implement handlers: log request → return dummy response / `Empty`

---

## Service: `chat-server` (Chat API)

### Features
- Chat management + sending messages over gRPC

### Proto
- [ ] `Create(usernames[]) -> (id)`
- [ ] `Delete(id) -> Empty`
- [ ] `SendMessage(from, text, timestamp) -> Empty`

### Implementation
- [ ] Implement handlers: log request → return dummy response / `Empty`

---

## Done criteria
- [ ] Both services compile and run
- [ ] `.proto` files fully describe the APIs above
- [ ] Stubs are generated successfully
- [ ] All RPCs are implemented (log + dummy response is OK)
