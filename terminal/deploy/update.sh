#!/usr/bin/env bash
# Run from /opt/hyr-remote-console/releases/<sha>-<run>-<attempt>/update.sh.
# Updates only the dedicated remote console service; never modifies admin SSH.
set -euo pipefail

RELEASE_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
APP_DIR="$(cd -- "$RELEASE_DIR/../.." && pwd)"
RELEASE_ID="$(basename -- "$RELEASE_DIR")"
[[ "$RELEASE_ID" =~ ^[0-9a-f]{40}-[0-9]+-[0-9]+$ ]] || { echo "Invalid release directory" >&2; exit 1; }
[[ "$(basename -- "$(dirname -- "$RELEASE_DIR")")" == releases ]] || exit 1
[[ -f "$APP_DIR/.env" ]] || { echo "Create $APP_DIR/.env before deploying" >&2; exit 1; }
[[ ! -e "$APP_DIR/current" || -L "$APP_DIR/current" ]] || { echo "current must be a symlink" >&2; exit 1; }
exec 9>"$APP_DIR/deploy.lock"
flock -n 9 || { echo "A console deployment is already running" >&2; exit 1; }

PREVIOUS="$(readlink -f "$APP_DIR/current" || true)"
compose() {
    docker compose --project-name hyr-remote-console --env-file "$APP_DIR/.env" -f "$1/docker-compose.yml" "${@:2}"
}
export CONSOLE_IMAGE="hyr-remote-console:$RELEASE_ID"
compose "$RELEASE_DIR" config --quiet
docker image load --input "$RELEASE_DIR/console-image.tar.gz"
docker image inspect "$CONSOLE_IMAGE" >/dev/null
# The loaded image is retained; remove only this release's redundant transfer archive.
rm -- "$RELEASE_DIR/console-image.tar.gz"

if compose "$RELEASE_DIR" up -d --no-build --pull never --no-deps --wait --wait-timeout 60 console; then
    ln -sfn "releases/$RELEASE_ID" "$APP_DIR/current"
    echo "Console healthy: $RELEASE_ID"
else
    compose "$RELEASE_DIR" logs --tail 30 console || true
    if [[ "$PREVIOUS" == "$APP_DIR/releases/"* && -f "$PREVIOUS/docker-compose.yml" ]]; then
        export CONSOLE_IMAGE="hyr-remote-console:$(basename -- "$PREVIOUS")"
        echo "Restoring previous console release"
        compose "$PREVIOUS" up -d --no-build --pull never --no-deps --wait --wait-timeout 60 console
    else
        compose "$RELEASE_DIR" stop console
    fi
    exit 1
fi
