# ==========================================
# STAGE 1: Shared Base Environment
# ==========================================
FROM golang:1.25.5-alpine AS base
WORKDIR /app

# Cache dependencies separately for faster builds
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire project source
COPY . .

# ==========================================
# STAGE 2: Build Binaries (Parallelizable)
# ==========================================
FROM base AS build-chat
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o grpc-chat-server ./cmd/chat/main.go

FROM base AS build-user
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o grpc-user-server ./cmd/user/main.go

# ==========================================
# STAGE 3: Final Chat Service Image
# ==========================================
FROM alpine:3.20 AS chat
WORKDIR /app
COPY --from=build-chat /app/grpc-chat-server .

# Run as non-root for security
USER 1000:1000
CMD ["./grpc-chat-server"]

# ==========================================
# STAGE 4: Final User Service Image
# ==========================================
FROM alpine:3.20 AS user
WORKDIR /app
COPY --from=build-user /app/grpc-user-server .

# Run as non-root for security
USER 1000:1000
CMD ["./grpc-user-server"]