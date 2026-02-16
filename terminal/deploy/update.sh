#!/bin/bash
# ─────────────────────────────────────────────────────────────
# Update Script — rebuild and restart the portfolio
# Usage: ./deploy/update.sh
# ─────────────────────────────────────────────────────────────
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

echo "Rebuilding portfolio container..."
docker compose build --no-cache portfolio

echo "Restarting portfolio container..."
docker compose up -d portfolio

echo ""
echo "Waiting for health check..."
sleep 5

STATUS=$(docker inspect --format='{{.State.Health.Status}}' terminal-portfolio 2>/dev/null || echo "unknown")

if [ "$STATUS" = "healthy" ]; then
    echo "Portfolio updated and healthy!"
else
    echo "Container status: $STATUS"
    echo "Check logs: docker compose logs portfolio --tail 20"
fi

echo ""
echo "Container Status:"
docker compose ps --format "table {{.Name}}\t{{.Status}}"
