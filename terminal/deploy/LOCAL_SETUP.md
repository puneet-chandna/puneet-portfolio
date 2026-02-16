# Local Deployment Guide — Terminal Portfolio

Run your terminal portfolio locally via Docker and expose it to the world through **Cloudflare Tunnel**.

## Prerequisites

| Requirement                | Check                                             |
| -------------------------- | ------------------------------------------------- |
| Docker + Docker Compose v2 | `docker compose version`                          |
| Domain on Cloudflare DNS   | `puneet.space` must be managed through Cloudflare |
| SSH host keys              | Already in `terminal/.ssh/`                       |
| API keys in `.env`         | Already in `terminal/.env`                        |

## Quick Start

### 1. One-Time Setup

```bash
cd terminal
./deploy/local-setup.sh
```

This will:

- Install `cloudflared` (if missing)
- Open a browser to authenticate with Cloudflare
- Create a named tunnel (`terminal-portfolio`)
- Configure DNS for `puneet.space` → your tunnel
- Save the tunnel token to `.env`
- Build and launch both containers

### 2. Verify

```bash
# From anywhere:
ssh puneet.space
```

## Updating After Code Changes

```bash
cd terminal
git pull                   # if needed
./deploy/update.sh         # rebuild + restart (tunnel stays up)
```

## Architecture

```
Internet → Cloudflare Edge (puneet.space:22)
               ↓ encrypted tunnel (QUIC)
         cloudflared container
               ↓ Docker bridge network
         portfolio container (port 22)
```

**WiFi changes?** The tunnel reconnects automatically — no static IP or port-forwarding needed.

## Common Commands

```bash
# View logs
docker compose logs -f                  # all services
docker compose logs -f portfolio        # just the portfolio
docker compose logs -f cloudflared      # just the tunnel

# Restart everything
docker compose restart

# Stop everything
docker compose down

# Full rebuild (nuclear option)
docker compose down
docker compose build --no-cache
docker compose up -d
```

## Troubleshooting

| Problem                    | Fix                                                                                  |
| -------------------------- | ------------------------------------------------------------------------------------ |
| `ssh puneet.space` hangs   | Check tunnel: `docker compose logs cloudflared`                                      |
| Container keeps restarting | Check app logs: `docker compose logs portfolio`                                      |
| Tunnel token expired       | Re-run `./deploy/local-setup.sh`                                                     |
| Port 22 already in use     | Your system SSH is on 22. Add `SSH_PORT=2222` to `.env` and update the tunnel config |

## Files Overview

| File                    | Purpose                                       |
| ----------------------- | --------------------------------------------- |
| `Dockerfile`            | Multi-stage Go build → minimal Alpine runtime |
| `docker-compose.yml`    | Portfolio + Cloudflare Tunnel services        |
| `.dockerignore`         | Keeps build context clean                     |
| `deploy/local-setup.sh` | One-time Cloudflare Tunnel setup              |
| `deploy/update.sh`      | Quick rebuild after code changes              |
