#!/usr/bin/env bash
# ============================================================
# Install Docker Engine + Docker Compose on Ubuntu
#
# Usage:
#   sudo bash install-docker-ubuntu.sh
#
# What it does:
#   1. Removes conflicting packages (docker.io, podman, etc.)
#   2. Adds the official Docker apt repository
#   3. Installs docker-ce, docker-compose-plugin, buildx
#   4. Adds your user to the "docker" group (no sudo needed)
# ============================================================

set -euo pipefail
trap 'echo "[!] Failed at line $LINENO. See output above." >&2; exit 1' ERR

# Must run as root
if [[ "${EUID}" -ne 0 ]]; then
  echo "[!] Please run with sudo: sudo bash $0"
  exit 1
fi

# Detect the real user behind sudo (for docker group)
USER_TO_ADD="${SUDO_USER:-$(id -un)}"

# Verify we're on Ubuntu and get the codename
# shellcheck disable=SC1091
. /etc/os-release
if [[ "${ID:-}" != "ubuntu" ]]; then
  echo "[!] This script is for Ubuntu. Detected: ${ID:-unknown}"
  exit 1
fi
CODENAME="${UBUNTU_CODENAME:-${VERSION_CODENAME}}"

echo "[1/6] Removing conflicting packages..."
apt-get remove -y docker.io docker-compose docker-compose-v2 \
  docker-doc podman-docker containerd runc 2>/dev/null || true

echo "[2/6] Installing prerequisites..."
apt-get update -y
apt-get install -y ca-certificates curl

echo "[3/6] Adding official Docker apt repository..."
install -m 0755 -d /etc/apt/keyrings
curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
chmod a+r /etc/apt/keyrings/docker.asc

tee /etc/apt/sources.list.d/docker.sources >/dev/null <<EOF
Types: deb
URIs: https://download.docker.com/linux/ubuntu
Suites: ${CODENAME}
Components: stable
Signed-By: /etc/apt/keyrings/docker.asc
EOF

echo "[4/6] Installing Docker Engine + Compose..."
apt-get update -y
apt-get install -y docker-ce docker-ce-cli containerd.io \
  docker-buildx-plugin docker-compose-plugin

echo "[5/6] Enabling Docker service..."
systemctl enable --now docker

echo "[6/6] Adding '${USER_TO_ADD}' to docker group..."
usermod -aG docker "${USER_TO_ADD}"

echo ""
echo "Done! Verifying:"
docker --version
docker compose version
echo ""
echo "Log out and back in (or run: newgrp docker) for group changes to take effect."
