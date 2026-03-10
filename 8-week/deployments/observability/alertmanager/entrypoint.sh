#!/bin/sh
# Replaces __TELEGRAM_*__ placeholders in alertmanager.yml
# with real values from environment variables, then starts
# Alertmanager. This keeps secrets out of Git.

# exit on error (-e), fail on unset variables (-u)
set -eu

# fail if required env vars are missing
: "${ALERTMANAGER_TELEGRAM_BOT_TOKEN:?ALERTMANAGER_TELEGRAM_BOT_TOKEN is required}"
: "${ALERTMANAGER_TELEGRAM_CHAT_ID:?ALERTMANAGER_TELEGRAM_CHAT_ID is required}"

# escape special sed characters in credential values
escaped_bot_token=$(printf "%s" "${ALERTMANAGER_TELEGRAM_BOT_TOKEN}" | sed -e 's/[\/&|]/\\&/g')
escaped_chat_id=$(printf "%s" "${ALERTMANAGER_TELEGRAM_CHAT_ID}" | sed -e 's/[\/&|]/\\&/g')

# replace placeholders and write to /tmp (writable)
sed \
  -e "s|__TELEGRAM_BOT_TOKEN__|${escaped_bot_token}|g" \
  -e "s|__TELEGRAM_CHAT_ID__|${escaped_chat_id}|g" \
  /etc/alertmanager/alertmanager.yml > /tmp/alertmanager.yml

# start alertmanager (exec replaces shell for proper signal handling)
exec /bin/alertmanager --config.file=/tmp/alertmanager.yml --storage.path=/alertmanager
