#!/bin/sh

set -eu

: "${ALERTMANAGER_TELEGRAM_BOT_TOKEN:?ALERTMANAGER_TELEGRAM_BOT_TOKEN is required}"
: "${ALERTMANAGER_TELEGRAM_CHAT_ID:?ALERTMANAGER_TELEGRAM_CHAT_ID is required}"

escaped_bot_token=$(printf "%s" "${ALERTMANAGER_TELEGRAM_BOT_TOKEN}" | sed -e 's/[\/&|]/\\&/g')
escaped_chat_id=$(printf "%s" "${ALERTMANAGER_TELEGRAM_CHAT_ID}" | sed -e 's/[\/&|]/\\&/g')

sed \
  -e "s|__TELEGRAM_BOT_TOKEN__|${escaped_bot_token}|g" \
  -e "s|__TELEGRAM_CHAT_ID__|${escaped_chat_id}|g" \
  /etc/alertmanager/alertmanager.yml > /tmp/alertmanager.yml

exec /bin/alertmanager --config.file=/tmp/alertmanager.yml --storage.path=/alertmanager
