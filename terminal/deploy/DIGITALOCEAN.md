# DigitalOcean Docker deployment

The public endpoint is `ssh puneet.space`. The portfolio uses **host port 22 → container port 2222**. Administrative OpenSSH uses a separate verified port, normally **2222**, on the droplet IP or an admin hostname.

## Coexisting with other containers

The Compose project is explicitly `puneet-terminal`. It creates its own bridge network and persistent `hostkeys` volume. The container has no host networking, Docker socket, privileged mode, or shared application volumes. It is limited to 0.5 CPU, 256 MiB RAM, 64 processes, and 30 MiB of rotated container logs. GitHub builds the image, so builds do not compete with production apps on the droplet.

`update.sh` targets only this project's `portfolio` service. It never runs Compose `down`, global prune, OS upgrades, or Docker restarts. A failed health check triggers restoration of the previous portfolio image/config when available. Existing sessions may disconnect during an update; this is a single-container deployment, not zero downtime.

**Port 22 is shared by every hostname resolving to the same droplet IP.** After cutover, `ssh root@DROPLET_IP` on port 22 reaches the portfolio too. Every existing CI deploy, backup, SFTP client, monitoring check, and SSH alias that administers this droplet must use the new admin port. Containers on other host ports can keep running, but that is only confirmed after checking their health before and after cutover. CPU, disk, and network capacity are still shared; limits reduce risk, not eliminate it.

## 1. Inspect before changing the host

Required: Linux x86_64, working Docker Engine, Compose v2 with `up --wait`, `flock`, enough free disk/RAM, and an administrative user able to manage Docker and `/opt/terminal-portfolio`. Docker group membership effectively grants root access; use a dedicated trusted deployment key/user.

Collect current state over the existing admin connection:

```bash
uname -m
cat /etc/os-release
docker version
docker compose version
docker ps --format 'table {{.Names}}\t{{.Ports}}\t{{.Status}}'
docker stats --no-stream
free -h
df -h /
sudo ss -ltnp
sudo sshd -t
sudo sshd -T | grep -E '^(port|listenaddress|passwordauthentication|pubkeyauthentication|permitrootlogin) '
systemctl is-active ssh.socket ssh.service sshd.service
```

Record the existing apps' direct and public health checks. Check both the DigitalOcean cloud firewall and host firewall. If Docker or Compose is missing, inspect package provenance and proposed package changes first; do not reinstall or upgrade Docker on a shared production host just to get Compose.

## 2. Preserve admin SSH before freeing port 22

This migration is intentionally **not automated by this repository**:

1. Keep the current admin session open. Verify access to DigitalOcean's **Recovery Console** first. The ordinary Droplet Console uses SSH and needs its custom-port configuration updated too.
2. Confirm host port 2222 is free. Open it in the relevant firewalls for trusted admin/deploy sources, preserving every existing application rule. GitHub-hosted runners do not have a single fixed egress IP: a laptop-only allowlist will block CI. Set an approved runner-access policy before enabling deployment.
3. Configure OpenSSH to listen on both 22 and 2222 temporarily. Inspect existing `Port`, `ListenAddress`, included config files, and systemd socket activation before choosing commands. On socket-activated Ubuntu, changing `sshd_config` alone may not change the listener; the socket configuration must agree.
4. Validate `sshd -t`, then reload the appropriate service/socket using host-specific steps. In a **second independent terminal**, verify `ssh -p 2222 ADMIN_USER@DROPLET_IP`, sudo access, and the deploy user's access.
5. Update every existing admin SSH client and automation for this droplet. Confirm each deployment path works on 2222.
6. Remove only the admin listener on port 22, validate again, and check `ss -ltnp` to ensure 22 is actually free. Do not close the original session until the new connection is proven.

If something else must remain on port 22, stop here: it cannot share the same address/port with this TUI. Use a separate IP/droplet or a nonstandard public portfolio port instead.

Docker publishes ports through its own firewall rules and can bypass UFW's normal filtering. Use the DigitalOcean cloud firewall as an outer boundary; do not rewrite Docker's firewall chains. Preserve rules for existing apps.

## 3. Prepare the app directory and optional email

Create `/opt/terminal-portfolio/releases` owned by the trusted deploy user, without changing ownership of other app directories. Create `/opt/terminal-portfolio/.env` with mode `600` and the same owner.

Start with an unused localhost-only port for staging:

```dotenv
PORTFOLIO_BIND_IP=127.0.0.1
PORTFOLIO_PORT=22220
RESEND_API_KEY=
RESEND_FROM=
RESEND_TO=
```

Email is disabled unless all three Resend values are configured. Use a verified sender such as `Puneet Portfolio <portfolio@YOUR_VERIFIED_DOMAIN>` and your actual recipient address. Do not use `onboarding@resend.dev` as a production default. Resend accepting a request is not proof of inbox delivery; verify one explicitly authorized live message before claiming delivery works. The web portfolio's Web3Forms key is unrelated.

### Preserve an existing portfolio identity

The new named volume does **not** automatically import an old bind-mounted `terminal/.ssh/id_ed25519`. If this portfolio already served visitors, securely back up that **portfolio** key and copy it into `puneet-terminal_hostkeys` before the first container starts, preserving the filename `id_ed25519`, ownership `10001:10001`, and file mode `600`. Copy its `.pub` file too if present. Never substitute the droplet's administrative SSH host key.

Stop only the old portfolio process during migration. Inspect the target volume first; do not overwrite an existing key. Verify the old and new public-key fingerprints match before publishing the new container. Keep the old key/backup until verification succeeds. Starting with an empty volume creates a new identity and causes host-key warnings for returning visitors; do that only as an intentional key rotation. Neither this workflow nor the update script deletes or imports local `.ssh` files.

Keep `.env` only on the droplet. The workflow transfers the tested image, Compose file, and update script; it does not copy, print, or replace `.env`.

## 4. Configure GitHub Actions

Create the `production` environment and use its deployment protection rules as appropriate. The deploy user must already have Docker permission and write access to the app directory; the workflow does not elevate privileges or install packages.

| Setting | Type | Value |
| --- | --- | --- |
| `DROPLET_HOST` | Secret | Droplet IPv4 address or a stable admin hostname, without a username or port |
| `DROPLET_USER` | Secret | Trusted deploy username |
| `DROPLET_SSH_PRIVATE_KEY` | Secret | Dedicated admin deployment key, whose public key is authorized on the droplet |
| `DROPLET_KNOWN_HOSTS` | Secret | Verified OpenSSH known-hosts entry for `[DROPLET_HOST]:2222` (or your chosen admin port) |
| `DROPLET_SSH_PORT` | Variable | `2222`, or the verified alternative; port 22 is intentionally refused |
| `PORTFOLIO_DEPLOY_ENABLED` | **Repository variable** | `true` only after prerequisites are verified; unset/false runs CI without deployment |

Obtain the admin host-key fingerprint over an already trusted connection or Recovery Console (`ssh-keygen -lf /etc/ssh/ssh_host_ed25519_key.pub`). Collect the matching known-hosts line for the exact deployment address and port and compare fingerprints before saving it. An unverified `ssh-keyscan` result alone is not authentication. Admin host keys and portfolio host keys are different identities; never reuse the admin private host key in the container.

The workflow runs on relevant PRs, pushes to `main`, and manual dispatch. It checks modules, vet, race tests, and govulncheck, then builds and smoke-tests the container. Successful main runs can transfer a uniquely tagged image and deploy it. Failed tests/scans stop deployment. Old `EC2_*` secrets are no longer referenced; remove them separately if no other workflow uses them.

The update script locks this app's deployments, loads the image, runs `up --no-build --pull never --no-deps --wait` for `portfolio`, and switches `current` only after health succeeds. It restores the previous release on failure and exits nonzero. It keeps the host-key volume and runtime `.env`; rollback uses the current `.env`, not a saved copy of old secrets. Image transfer archives are deleted only after a successful load; previous images/releases remain for rollback. Periodically review and remove only obsolete **portfolio** releases/images; never run global Docker prune on this shared host.

## 5. Test staging, then publish

For the localhost-only staging port, open an admin tunnel from your machine:

```bash
ssh -p 2222 -L 22220:127.0.0.1:22220 ADMIN_USER@DROPLET_IP
# In another local terminal:
ssh -p 22220 localhost
```

Check navigation, resize, help, contact availability, and quitting. Restart/recreate only the portfolio service and confirm its SSH fingerprint stays unchanged. Check all other apps against the baseline.

After admin port migration is complete and port 22 is free, update the droplet's `.env` to:

```dotenv
PORTFOLIO_BIND_IP=0.0.0.0
PORTFOLIO_PORT=22
```

Keep any configured Resend values. Re-run the workflow to apply this binding. Permit inbound TCP 22 for public visitors in the cloud firewall; keep the admin port's separate rules. Do not alter existing HTTP/HTTPS or application port rules.

In Cloudflare DNS:

| Type | Name | Content | Proxy |
| --- | --- | --- | --- |
| A | `@` | Droplet public IPv4 | **DNS only** |

Replace the previous apex tunnel/proxy record after checking its use. Preserve unrelated DNS records, including mail and verification records. Remove a stale apex AAAA record unless the droplet's IPv6 listener/firewall path is configured and tested. Do not publish an untested IPv6 address. This guide's default binding is IPv4.

Normal Cloudflare Tunnel SSH requires client-side setup, so it is not used for this public plain-SSH endpoint. DNS does not select a different IP for web and SSH traffic: if `puneet.space` also needs a website, that HTTP/HTTPS service or redirect must be served on the same destination IP. The separate `puneetchandna.com` website is unchanged.

Finally, from outside the droplet:

```bash
ssh puneet.space
ssh -p 2222 ADMIN_USER@DROPLET_IP
```

Check DNS A/AAAA records, both SSH identities, container health after recreation, and every existing app's health. A first-time SSH host-key confirmation is normal; compare the expected fingerprint. If replacing an old portfolio host key, verify the new identity before removing an old known-hosts entry.

## Operations

`/opt/terminal-portfolio/current` points to the last healthy release. Use `docker compose --project-name puneet-terminal --env-file /opt/terminal-portfolio/.env -f /opt/terminal-portfolio/current/docker-compose.yml ps` or `logs --tail 50 portfolio` for inspection. Before any `up`, set `PORTFOLIO_IMAGE=puneet-terminal:RELEASE_ID` for the intended release; without it the local-development image default is used.

To roll back manually, pick a retained release ID, set that `PORTFOLIO_IMAGE`, and run the same scoped `up -d --no-build --pull never --no-deps --wait --wait-timeout 60 portfolio` against its Compose file. Update `current` only after success. Never remove `puneet-terminal_hostkeys`, use `down -v`, or delete that volume while cleaning old releases. Back up its private host key securely, outside this repository.

References: [Compose project isolation](https://docs.docker.com/compose/how-tos/project-name/), [Docker firewall behavior](https://docs.docker.com/engine/network/packet-filtering-firewalls/), [DigitalOcean console and custom SSH ports](https://docs.digitalocean.com/products/droplets/how-to/connect-with-console/), [Cloudflare DNS-only SSH records](https://developers.cloudflare.com/dns/proxy-status/use-cases/).
