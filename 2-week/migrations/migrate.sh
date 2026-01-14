#!/bin/sh
set -eu

MIGRATIONS_DIR="${MIGRATIONS_DIR:-/migrations}"
ENV_FILE="${ENV_FILE:-../local.env}"

# Load env from file (only if POSTGRES_HOST isn't set)
if [ -z "${POSTGRES_HOST:-}" ] && [ -f "$ENV_FILE" ]; then
  echo "🧾 Loading env from $ENV_FILE"
  set -a
  . "$ENV_FILE"
  set +a
fi

# Validate required vars
: "${POSTGRES_HOST:?🛑 POSTGRES_HOST is not set}"
: "${POSTGRES_PORT:?🛑 POSTGRES_PORT is not set}"
: "${POSTGRES_USER:?🛑 POSTGRES_USER is not set}"
: "${POSTGRES_PASSWORD:?🛑 POSTGRES_PASSWORD is not set}"
: "${POSTGRES_DB:?🛑 POSTGRES_DB is not set}"
POSTGRES_SSLMODE="${POSTGRES_SSLMODE:-disable}"

echo "⏳ Waiting for Postgres at ${POSTGRES_HOST}:${POSTGRES_PORT}..."
until nc -z -w 2 "$POSTGRES_HOST" "$POSTGRES_PORT" >/dev/null 2>&1; do
  echo "💤 Not ready yet, retrying in 2s..."
  sleep 2
done
echo "✅ Postgres is up"

DSN="host=$POSTGRES_HOST port=$POSTGRES_PORT user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB sslmode=$POSTGRES_SSLMODE"

echo "🧩 Checking migration status..."
STATUS="$(goose -dir "$MIGRATIONS_DIR" postgres "$DSN" status 2>&1)"
echo "$STATUS"

if echo "$STATUS" | grep -qE '(^|[[:space:]])Pending([[:space:]]+--|[[:space:]]|$)'; then
  echo "🚀 Pending migrations found — applying..."
  goose -dir "$MIGRATIONS_DIR" postgres "$DSN" up
  echo "🎉 Migrations applied"
else
  echo "✅ No pending migrations — nothing to do"
fi
