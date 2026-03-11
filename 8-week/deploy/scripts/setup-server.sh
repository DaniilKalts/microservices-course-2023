#!/usr/bin/env bash
# ============================================================
# setup-server.sh — One-shot Ubuntu server provisioning
#
# Run on a fresh Ubuntu server:
#   sudo bash scripts/setup-server.sh
#
# What it does:
#   1. Configures UFW firewall (SSH + HTTP + HTTPS only)
#   2. Hardens SSH (disables root & password login)
#   3. Installs Go 1.25.x
#   4. Installs Task (task runner)
#   5. Installs Docker + Docker Compose (via install-docker-ubuntu.sh)
#   6. Creates .env from .env.example with secure random secrets
#   7. Generates JWT keys and TLS certs via Taskfile
#
# After this script finishes, you can run:
#   cd /path/to/project && docker compose up --build
# ============================================================

set -euo pipefail
trap 'echo "[!] Failed at line $LINENO. See output above." >&2; exit 1' ERR

# ── Must run as root ────────────────────────────────────────
if [[ "${EUID}" -ne 0 ]]; then
  echo "[!] Please run with sudo: sudo bash $0"
  exit 1
fi

REAL_USER="${SUDO_USER:-$(id -un)}"
REAL_HOME=$(eval echo "~${REAL_USER}")
SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

GO_VERSION="1.25.7"
TASK_VERSION="v3.49.0"

# ── Allowed SSH sources ─────────────────────────────────────
# Add your IPs/CIDRs here. Empty = allow SSH from anywhere.
ALLOWED_SSH_SOURCES=(
  # "203.0.113.10"
  # "198.51.100.0/24"
)

echo "============================================"
echo " Server setup for: ${PROJECT_DIR}"
echo " User: ${REAL_USER}"
echo "============================================"

# ── 1. Firewall (UFW) ──────────────────────────────────────
echo ""
echo "[1/7] Configuring firewall..."

apt-get update -y
apt-get install -y ufw

# Reset to clean state — deny all incoming, allow all outgoing.
ufw --force reset
ufw default deny incoming
ufw default allow outgoing

# SSH — restrict to specific IPs if configured, otherwise allow from anywhere.
if [[ ${#ALLOWED_SSH_SOURCES[@]} -gt 0 ]]; then
  for src in "${ALLOWED_SSH_SOURCES[@]}"; do
    ufw allow from "${src}" to any port 22 proto tcp comment "SSH from ${src}"
  done
else
  ufw allow 22/tcp comment "SSH"
fi

# HTTP + HTTPS — open to the world (app + Grafana + reverse proxy).
ufw allow 80/tcp comment "HTTP"
ufw allow 443/tcp comment "HTTPS"

# Enable firewall (--force skips the interactive prompt).
ufw --force enable
ufw status verbose

# ── 2. SSH hardening ────────────────────────────────────────
echo ""
echo "[2/7] Hardening SSH..."

SSHD_CONFIG="/etc/ssh/sshd_config"

# Disable root login and password auth (use SSH keys instead).
sed -i 's/^#\?PermitRootLogin.*/PermitRootLogin no/' "${SSHD_CONFIG}"
sed -i 's/^#\?PasswordAuthentication.*/PasswordAuthentication no/' "${SSHD_CONFIG}"

# Restart SSH to apply changes.
systemctl restart sshd

echo "    Root login: disabled"
echo "    Password auth: disabled (SSH keys only)"

# ── 3. Install Go ──────────────────────────────────────────
echo ""
echo "[3/7] Installing Go ${GO_VERSION}..."

if command -v go &>/dev/null && go version | grep -q "go${GO_VERSION}"; then
  echo "    Already installed: $(go version)"
else
  curl -fsSL "https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz" -o /tmp/go.tar.gz
  rm -rf /usr/local/go
  tar -C /usr/local -xzf /tmp/go.tar.gz
  rm /tmp/go.tar.gz

  # Add Go to PATH for the real user (idempotent).
  GO_ENV_LINE='export PATH=$PATH:/usr/local/go/bin:$(go env GOPATH)/bin'
  PROFILE="${REAL_HOME}/.profile"
  if ! grep -qF '/usr/local/go/bin' "${PROFILE}" 2>/dev/null; then
    echo "${GO_ENV_LINE}" >> "${PROFILE}"
  fi

  # Make Go available for the rest of this script.
  export PATH=$PATH:/usr/local/go/bin
  echo "    Installed: $(go version)"
fi

# ── 4. Install Task ────────────────────────────────────────
echo ""
echo "[4/7] Installing Task ${TASK_VERSION}..."

if command -v task &>/dev/null; then
  echo "    Already installed: $(task --version)"
else
  sh -c "$(curl -fsSL https://taskfile.dev/install.sh)" -- -d -b /usr/local/bin "${TASK_VERSION}"
  echo "    Installed: $(task --version)"
fi

# ── 5. Install Docker ──────────────────────────────────────
echo ""
echo "[5/7] Installing Docker..."

if command -v docker &>/dev/null; then
  echo "    Already installed: $(docker --version)"
else
  bash "${SCRIPT_DIR}/install-docker-ubuntu.sh"
fi

# ── 6. Create .env ──────────────────────────────────────────
echo ""
echo "[6/7] Creating .env..."

ENV_FILE="${PROJECT_DIR}/.env"
ENV_EXAMPLE="${PROJECT_DIR}/.env.example"

if [[ -f "${ENV_FILE}" ]]; then
  echo "    .env already exists — skipping (delete it to regenerate)"
else
  if [[ ! -f "${ENV_EXAMPLE}" ]]; then
    echo "[!] .env.example not found at ${ENV_EXAMPLE}"
    exit 1
  fi

  cp "${ENV_EXAMPLE}" "${ENV_FILE}"

  # Generate a secure random Postgres password.
  PG_PASS=$(openssl rand -base64 24 | tr -d '/+=')
  sed -i "s|^POSTGRES_PASSWORD=.*|POSTGRES_PASSWORD=${PG_PASS}|" "${ENV_FILE}"

  # Set Postgres host to "postgres" (docker-compose service name).
  sed -i "s|^POSTGRES_HOST=.*|POSTGRES_HOST=postgres|" "${ENV_FILE}"

  # Generate a secure Grafana admin password.
  GF_PASS=$(openssl rand -base64 16 | tr -d '/+=')
  sed -i "s|^GRAFANA_ADMIN_PASSWORD=.*|GRAFANA_ADMIN_PASSWORD=${GF_PASS}|" "${ENV_FILE}"

  # Restrict .env to owner-only access.
  chown "${REAL_USER}:${REAL_USER}" "${ENV_FILE}"
  chmod 600 "${ENV_FILE}"

  echo "    Created .env with random secrets"
  echo "    Postgres password: ${PG_PASS}"
  echo "    Grafana password:  ${GF_PASS}"
  echo ""
  echo "    IMPORTANT: Update these values in .env before going to production:"
  echo "      - ALERTMANAGER_TELEGRAM_BOT_TOKEN"
  echo "      - ALERTMANAGER_TELEGRAM_CHAT_ID"
  echo "      - TLS_CERTBOT_DOMAIN"
  echo "      - TLS_CERTBOT_EMAIL"
fi

# ── 7. Generate build artifacts ─────────────────────────────
echo ""
echo "[7/7] Generating JWT keys and TLS certs..."

# Run as the real user so file ownership is correct.
cd "${PROJECT_DIR}"
sudo -u "${REAL_USER}" task jwt:generate 2>/dev/null || echo "    (jwt:generate not available — skip)"
sudo -u "${REAL_USER}" task tls:generate 2>/dev/null || echo "    (tls:generate not available — skip)"

echo ""
echo "============================================"
echo " Setup complete!"
echo "============================================"
echo ""
echo " Next steps:"
echo "   1. Log out and back in (for docker group)"
echo "   2. Edit .env with your Telegram bot token"
echo "      and real domain/email for TLS"
echo "   3. Run the app:"
echo "      cd ${PROJECT_DIR}"
echo "      docker compose up -d"
echo ""
