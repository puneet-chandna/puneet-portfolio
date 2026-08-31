#!/usr/bin/env bash
# Run from /opt/terminal-portfolio/releases/<sha>-<run>-<attempt>/update.sh.
# Updates only the dedicated portfolio Compose service; never modifies admin SSH.
set -euo pipefail

RELEASE_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
APP_DIR="$(cd -- "$RELEASE_DIR/../.." && pwd)"
RELEASE_ID="$(basename -- "$RELEASE_DIR")"
[[ "$RELEASE_ID" =~ ^[0-9a-f]{40}-[0-9]+-[0-9]+$ ]] || { echo "Invalid release directory" >&2; exit 1; }
[[ "$(basename -- "$(dirname -- "$RELEASE_DIR")")" == releases ]] || exit 1
[[ -f "$APP_DIR/.env" ]] || { echo "Create $APP_DIR/.env before deploying" >&2; exit 1; }
[[ ! -e "$APP_DIR/current" || -L "$APP_DIR/current" ]] || { echo "current must be a symlink" >&2; exit 1; }
exec 9>"$APP_DIR/deploy.lock"
flock -n 9 || { echo "A portfolio deployment is already running" >&2; exit 1; }

PREVIOUS="$(readlink -f "$APP_DIR/current" || true)"
compose() {
    docker compose --project-name puneet-terminal --env-file "$APP_DIR/.env" -f "$1/docker-compose.yml" "${@:2}"
}
export PORTFOLIO_IMAGE="puneet-terminal:$RELEASE_ID"
compose "$RELEASE_DIR" config --quiet
docker image load --input "$RELEASE_DIR/portfolio-image.tar.gz"
docker image inspect "$PORTFOLIO_IMAGE" >/dev/null
# The loaded image is retained; remove only this release's redundant transfer archive.
rm -- "$RELEASE_DIR/portfolio-image.tar.gz"

if compose "$RELEASE_DIR" up -d --no-build --pull never --no-deps --wait --wait-timeout 60 portfolio; then
    ln -sfn "releases/$RELEASE_ID" "$APP_DIR/current"
    echo "Portfolio healthy: $RELEASE_ID"
else
    compose "$RELEASE_DIR" logs --tail 30 portfolio || true
    if [[ "$PREVIOUS" == "$APP_DIR/releases/"* && -f "$PREVIOUS/docker-compose.yml" ]]; then
        export PORTFOLIO_IMAGE="puneet-terminal:$(basename -- "$PREVIOUS")"
        echo "Restoring previous portfolio release"
        compose "$PREVIOUS" up -d --no-build --pull never --no-deps --wait --wait-timeout 60 portfolio
    else
        compose "$RELEASE_DIR" stop portfolio
    fi
    exit 1
fi
