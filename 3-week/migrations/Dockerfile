# ==========================================
# STAGE 1: Build Goose Migration Tool
# ==========================================
FROM golang:1.25.6-alpine AS builder

# Install goose migration tool from source
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# ==========================================
# STAGE 2: Final Migrator Image
# ==========================================
FROM alpine:3.20 AS migrator
WORKDIR /app

# Install netcat for database connection checking
RUN apk add --no-cache netcat-openbsd

# Copy goose binary from builder stage
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Copy the migration script and ensure it is executable
COPY migrations/migrate.sh .
RUN chmod +x ./migrate.sh

# Run the migration script on container start
ENTRYPOINT ["./migrate.sh"]