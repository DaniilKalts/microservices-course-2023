# STAGE 1: Build
# ==========================================
FROM golang:1.25.7-alpine AS builder
WORKDIR /app
RUN CGO_ENABLED=0 go install github.com/pressly/goose/v3/cmd/goose@v3.27.0

# === Run migrations ===
FROM alpine:3.20 AS migrator
WORKDIR /app

RUN apk add --no-cache postgresql-client

COPY --from=builder /go/bin/goose /usr/local/bin/goose
COPY deploy/migrations/migrate.sh .
RUN chmod +x ./migrate.sh

ENTRYPOINT ["./migrate.sh"]
