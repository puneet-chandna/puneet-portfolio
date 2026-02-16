#!/bin/bash
# ─────────────────────────────────────────────────────────────
# Local Setup Script for Terminal Portfolio
# Sets up Docker + Cloudflare Tunnel for puneet.space
# ─────────────────────────────────────────────────────────────
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"
ENV_FILE="$PROJECT_DIR/.env"
DOMAIN="puneet.space"
TUNNEL_NAME="terminal-portfolio"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m'

info()  { echo -e "${CYAN}[INFO]${NC} $1"; }
ok()    { echo -e "${GREEN}[  OK]${NC} $1"; }
warn()  { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[FAIL]${NC} $1"; exit 1; }

# ── Step 1: Check prerequisites ──────────────────────────────
info "Checking prerequisites..."

command -v docker >/dev/null 2>&1 || error "Docker is not installed. Install it first: https://docs.docker.com/get-docker/"
docker compose version >/dev/null 2>&1 || error "Docker Compose (v2) is not installed."
systemctl is-active docker >/dev/null 2>&1 || error "Docker daemon is not running. Start it with: sudo systemctl start docker"
ok "Docker is ready"

# ── Step 2: Install cloudflared if missing ───────────────────
if ! command -v cloudflared >/dev/null 2>&1; then
    info "Installing cloudflared..."
    
    # Detect architecture
    ARCH=$(uname -m)
    case "$ARCH" in
        x86_64)  CF_ARCH="amd64" ;;
        aarch64) CF_ARCH="arm64" ;;
        armv7l)  CF_ARCH="arm"   ;;
        *)       error "Unsupported architecture: $ARCH" ;;
    esac

    # Download and install
    CF_URL="https://github.com/cloudflare/cloudflared/releases/latest/download/cloudflared-linux-${CF_ARCH}.deb"
    TEMP_DEB=$(mktemp /tmp/cloudflared-XXXXXX.deb)
    
    curl -fsSL "$CF_URL" -o "$TEMP_DEB"
    sudo dpkg -i "$TEMP_DEB"
    rm -f "$TEMP_DEB"
    
    ok "cloudflared installed"
else
    ok "cloudflared already installed ($(cloudflared --version 2>&1 | head -1))"
fi

# ── Step 3: Authenticate with Cloudflare ─────────────────────
info "Checking Cloudflare authentication..."

if [ ! -f "$HOME/.cloudflared/cert.pem" ]; then
    warn "You need to log in to Cloudflare. A browser window will open."
    echo ""
    echo "  → Select the zone for: ${CYAN}${DOMAIN}${NC}"
    echo ""
    cloudflared tunnel login
    ok "Cloudflare authentication complete"
else
    ok "Already authenticated with Cloudflare"
fi

# ── Step 4: Create tunnel (or reuse existing) ────────────────
info "Setting up Cloudflare Tunnel: ${TUNNEL_NAME}..."

EXISTING_TUNNEL=$(cloudflared tunnel list -o json 2>/dev/null | python3 -c "
import json, sys
tunnels = json.load(sys.stdin)
for t in tunnels:
    if t['name'] == '${TUNNEL_NAME}':
        print(t['id'])
        break
" 2>/dev/null || echo "")

if [ -n "$EXISTING_TUNNEL" ]; then
    TUNNEL_ID="$EXISTING_TUNNEL"
    ok "Tunnel already exists: ${TUNNEL_ID}"
else
    info "Creating new tunnel..."
    TUNNEL_OUTPUT=$(cloudflared tunnel create "$TUNNEL_NAME" 2>&1)
    TUNNEL_ID=$(echo "$TUNNEL_OUTPUT" | grep -oP '[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}' | head -1)
    
    if [ -z "$TUNNEL_ID" ]; then
        error "Failed to create tunnel. Output: $TUNNEL_OUTPUT"
    fi
    ok "Tunnel created: ${TUNNEL_ID}"
fi

# ── Step 5: Get tunnel token ─────────────────────────────────
info "Retrieving tunnel token..."

TUNNEL_TOKEN=$(cloudflared tunnel token "$TUNNEL_NAME" 2>/dev/null)

if [ -z "$TUNNEL_TOKEN" ]; then
    error "Failed to get tunnel token"
fi
ok "Tunnel token retrieved"

# ── Step 6: Configure DNS ────────────────────────────────────
info "Configuring DNS route: ${DOMAIN} → tunnel..."

# Route the domain to the tunnel (creates CNAME record)
cloudflared tunnel route dns "$TUNNEL_NAME" "$DOMAIN" 2>/dev/null || {
    warn "DNS route may already exist, continuing..."
}
ok "DNS configured for ${DOMAIN}"

# ── Step 7: Create tunnel config ─────────────────────────────
info "Writing tunnel configuration..."

mkdir -p "$HOME/.cloudflared"
cat > "$HOME/.cloudflared/config.yml" <<EOF
tunnel: ${TUNNEL_ID}
credentials-file: /etc/cloudflared/${TUNNEL_ID}.json

ingress:
  - hostname: ${DOMAIN}
    service: ssh://portfolio:22
  - service: http_status:404
EOF

ok "Tunnel config written to ~/.cloudflared/config.yml (container path)"

# ── Step 8: Update .env with tunnel token ────────────────────
info "Updating .env file..."

# Ensure .env exists
touch "$ENV_FILE"

# Update or add CLOUDFLARE_TUNNEL_TOKEN
if grep -q "^CLOUDFLARE_TUNNEL_TOKEN=" "$ENV_FILE" 2>/dev/null; then
    sed -i "s|^CLOUDFLARE_TUNNEL_TOKEN=.*|CLOUDFLARE_TUNNEL_TOKEN=${TUNNEL_TOKEN}|" "$ENV_FILE"
else
    echo "CLOUDFLARE_TUNNEL_TOKEN=${TUNNEL_TOKEN}" >> "$ENV_FILE"
fi

# Ensure SSH_PORT is set
if ! grep -q "^SSH_PORT=" "$ENV_FILE" 2>/dev/null; then
    echo "SSH_PORT=22" >> "$ENV_FILE"
fi

ok ".env updated with tunnel token"

# ── Step 9: Build and launch ─────────────────────────────────
info "Building and launching containers..."

cd "$PROJECT_DIR"
docker compose build
docker compose up -d

ok "Containers are running!"

# ── Step 10: Verify ──────────────────────────────────────────
echo ""
echo "══════════════════════════════════════════════════════════"
echo ""
info "Waiting for services to stabilize..."
sleep 8

echo ""
echo "  Container Status:"
docker compose ps --format "table {{.Name}}\t{{.Status}}"

echo ""
echo "Tunnel Status:"
echo "     Dashboard: https://one.dash.cloudflare.com → Networks → Tunnels"
echo ""
echo " Test your portfolio:"
echo -e "     ${CYAN}ssh ${DOMAIN}${NC}"
echo ""
echo "  Update after code changes:"
echo -e "     ${CYAN}./deploy/update.sh${NC}"
echo ""
echo "══════════════════════════════════════════════════════════"
