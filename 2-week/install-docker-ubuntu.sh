#!/usr/bin/env bash
set -euo pipefail

# Docker Engine + Docker Compose v2 (docker compose) install script for Ubuntu.
# - Uses official Docker apt repo
# - Installs: docker-ce, docker-ce-cli, containerd.io, docker-buildx-plugin, docker-compose-plugin
# - Optional: adds a user to the "docker" group so you can run docker without sudo

# ---- Settings (override via env) ----
: "${DOCKER_ADD_USER:=1}"          # 1 = add user to docker group, 0 = skip
: "${DOCKER_USER:=}"              # if empty: uses SUDO_USER or current user
: "${RUN_HELLO_WORLD:=0}"         # 1 = run "docker run hello-world" at the end

# ---- Helpers ----
log() { echo -e "\n[+] $*\n"; }
die() { echo "[!] $*" >&2; exit 1; }

need_cmd() {
  command -v "$1" >/dev/null 2>&1 || die "Missing required command: $1"
}

# Run as root or via sudo
if [[ "${EUID}" -ne 0 ]]; then
  if command -v sudo >/dev/null 2>&1; then
    SUDO="sudo"
  else
    die "Please run as root, or install sudo and re-run."
  fi
else
  SUDO=""
fi

# Determine which user to add to docker group
if [[ -z "${DOCKER_USER}" ]]; then
  if [[ -n "${SUDO_USER:-}" && "${SUDO_USER}" != "root" ]]; then
    DOCKER_USER="${SUDO_USER}"
  else
    DOCKER_USER="$(id -un)"
  fi
fi

# Verify Ubuntu
if [[ -r /etc/os-release ]]; then
  # shellcheck disable=SC1091
  . /etc/os-release
else
  die "Cannot read /etc/os-release; are you sure this is Ubuntu?"
fi

if [[ "${ID:-}" != "ubuntu" ]]; then
  die "This script is for Ubuntu. Detected: ID=${ID:-unknown}"
fi

UBU_CODENAME="${UBUNTU_CODENAME:-${VERSION_CODENAME:-}}"
if [[ -z "${UBU_CODENAME}" ]]; then
  die "Cannot detect Ubuntu codename (UBUNTU_CODENAME/VERSION_CODENAME)."
fi

trap 'die "Failed at line $LINENO. See output above."' ERR

log "Updating apt index"
$SUDO apt-get update -y

log "Removing conflicting/unofficial packages if present"
conflicts=(docker.io docker-compose docker-compose-v2 docker-doc podman-docker containerd runc)
to_remove=()
for p in "${conflicts[@]}"; do
  if dpkg -s "$p" >/dev/null 2>&1; then
    to_remove+=("$p")
  fi
done
if ((${#to_remove[@]})); then
  $SUDO apt-get remove -y "${to_remove[@]}"
else
  echo "[i] No conflicting packages found."
fi

log "Installing prerequisites"
$SUDO apt-get install -y ca-certificates curl

log "Setting up Docker apt repository (official)"
$SUDO install -m 0755 -d /etc/apt/keyrings
$SUDO curl -fsSL https://download.docker.com/linux/ubuntu/gpg -o /etc/apt/keyrings/docker.asc
$SUDO chmod a+r /etc/apt/keyrings/docker.asc

# Create /etc/apt/sources.list.d/docker.sources (new-style deb822 sources file)
$SUDO tee /etc/apt/sources.list.d/docker.sources >/dev/null <<EOF
Types: deb
URIs: https://download.docker.com/linux/ubuntu
Suites: ${UBU_CODENAME}
Components: stable
Signed-By: /etc/apt/keyrings/docker.asc
EOF

log "Updating apt index after adding Docker repo"
$SUDO apt-get update -y

log "Installing Docker Engine + Compose v2 plugin + Buildx"
$SUDO apt-get install -y docker-ce docker-ce-cli containerd.io docker-buildx-plugin docker-compose-plugin

log "Enabling and starting Docker service"
$SUDO systemctl enable --now docker

if [[ "${DOCKER_ADD_USER}" == "1" ]]; then
  log "Adding user '${DOCKER_USER}' to docker group (so you can run docker without sudo)"
  # group 'docker' is usually created by the package, but ensure it exists
  if ! getent group docker >/dev/null 2>&1; then
    $SUDO groupadd docker
  fi
  $SUDO usermod -aG docker "${DOCKER_USER}"
  echo "[i] NOTE: You must log out and log back in (or run: newgrp docker) for group changes to take effect."
else
  echo "[i] Skipping docker group setup (DOCKER_ADD_USER=0)."
fi

log "Verifying installation"
docker --version
docker compose version

if [[ "${RUN_HELLO_WORLD}" == "1" ]]; then
  log "Running hello-world test"
  $SUDO docker run --rm hello-world
fi

log "Done ✅"
echo "[i] Next steps:"
echo "    - Reconnect SSH session (or: newgrp docker) if you added user to docker group."
echo "    - Test: docker ps"
echo "    - Compose: docker compose up -d"
