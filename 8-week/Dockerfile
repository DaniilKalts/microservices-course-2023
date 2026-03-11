# === Build ===
FROM golang:1.25.7-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o api-server ./cmd/main.go

# === Run ===
FROM alpine:3.20
WORKDIR /app

COPY --from=builder --chown=1000:1000 /app/api-server ./api-server
COPY --from=builder --chown=1000:1000 /app/api/gen/openapi ./api/gen/openapi

USER 1000:1000
ENTRYPOINT ["./api-server"]
