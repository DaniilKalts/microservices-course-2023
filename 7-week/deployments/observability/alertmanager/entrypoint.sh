#!/bin/sh
# ============================================================
# entrypoint.sh — Inject Telegram credentials into alertmanager.yml
#
# Why not put secrets directly in the YAML?
#   alertmanager.yml is committed to Git with placeholders.
#   This script swaps them for real values from environment
#   variables (loaded from .env by docker-compose) at container
#   startup, so secrets never touch version control.
#
# Flow:
#   .env → docker-compose → env vars → this script → sed → /tmp/alertmanager.yml → alertmanager starts
# ============================================================

set -eu   # exit on any error (-e), treat unset vars as errors (-u)

# Fail fast if required env vars are missing.
: "${ALERTMANAGER_TELEGRAM_BOT_TOKEN:?ALERTMANAGER_TELEGRAM_BOT_TOKEN is required}"
: "${ALERTMANAGER_TELEGRAM_CHAT_ID:?ALERTMANAGER_TELEGRAM_CHAT_ID is required}"

# Escape special sed characters (/ & |) in credential values
# so they don't break the substitution command.
escaped_bot_token=$(printf "%s" "${ALERTMANAGER_TELEGRAM_BOT_TOKEN}" | sed -e 's/[\/&|]/\\&/g')
escaped_chat_id=$(printf "%s" "${ALERTMANAGER_TELEGRAM_CHAT_ID}" | sed -e 's/[\/&|]/\\&/g')

# Replace placeholders in the read-only mounted config and write
# the result to /tmp (writable). Uses | as sed delimiter to avoid
# conflicts with / in tokens.
sed \
  -e "s|__TELEGRAM_BOT_TOKEN__|${escaped_bot_token}|g" \
  -e "s|__TELEGRAM_CHAT_ID__|${escaped_chat_id}|g" \
  /etc/alertmanager/alertmanager.yml > /tmp/alertmanager.yml

# exec replaces this shell with alertmanager (PID 1 for proper signal handling).
exec /bin/alertmanager --config.file=/tmp/alertmanager.yml --storage.path=/alertmanager
