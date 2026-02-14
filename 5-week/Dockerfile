# ==========================================
# STAGE 1: Build
# ==========================================
FROM golang:1.25.7-alpine AS builder
WORKDIR /app

# Cache dependencies separately for faster builds
COPY go.mod go.sum ./
RUN go mod download

# # Copy the entire project source
COPY . .

# FROM base AS build-chat
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o api ./cmd/main.go


# ==========================================
# STAGE 2: Run
# ==========================================
FROM alpine:3.20 AS chat
WORKDIR /app
COPY --from=builder /app/api .

# Run as non-root for security
USER 1000:1000
ENTRYPOINT ["./api"]