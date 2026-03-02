# Dockerize Urja for Production Deployment

**Date:** 2026-03-02
**Goal:** Production-ready Docker setup for Urja (Go API + Next.js web) on a single VPS, with a shared Traefik reverse proxy that serves 10+ apps, plus GitHub Actions CI/CD.

## Network Architecture

```
VPS
├── traefik/                      ← Shared infra (one per VPS)
│   └── docker-compose.yml        ← Traefik v3 + external "proxy" network
│
├── urja/                         ← This app
│   └── docker-compose.prod.yml   ← api, web, postgres, redis
│
├── app2/                         ← Future apps join same "proxy" network
└── ...
```

One external Docker network `proxy` created by Traefik. Every app joins it. Each app keeps its own internal network for DB/Redis isolation.

```
Internet → :80/:443 → Traefik
                        ├─ Host(api.domain.com)  → urja-api:8080
                        ├─ Host(domain.com)      → urja-web:3000
                        ├─ Host(otherapp.com)     → other-app:PORT
                        └─ ...

urja internal network:
  api ↔ postgres
  api ↔ redis
  (postgres/redis NOT on proxy network)
```

## Containers

### API (`Dockerfile` — modify existing)

- Keep existing multi-stage build (golang:1.22-alpine builder, alpine:3.19 runtime)
- Add `golang-migrate` CLI to runtime image
- Entrypoint script: run `migrate -path ./migrations -database $DB_URL up`, then exec `./urja-api`
- Migrations bundled in image (already done)

### Web (`web/Dockerfile` — new)

- Multi-stage: `node:20-alpine` builder → `node:20-alpine` runner
- Enable Next.js `output: "standalone"` in `next.config.mjs`
- Build-time arg: `NEXT_PUBLIC_API_URL` for API endpoint
- Final image ~100MB (standalone output, no node_modules)

### Postgres

- `postgres:16-alpine`, persistent named volume
- Internal network only, no external port exposure
- Healthcheck: `pg_isready -U urja`

### Redis

- `redis:7-alpine`, persistent named volume
- Internal network only
- Healthcheck: `redis-cli ping`

## Production Compose (`docker-compose.prod.yml`)

```yaml
services:
  api:
    image: ghcr.io/OWNER/urja-api:latest
    env_file: .env.prod
    depends_on:
      postgres: { condition: service_healthy }
      redis: { condition: service_healthy }
    networks: [proxy, internal]
    labels:
      - traefik.enable=true
      - traefik.http.routers.urja-api.rule=Host(`api.${DOMAIN}`)
      - traefik.http.routers.urja-api.tls.certresolver=letsencrypt
      - traefik.http.services.urja-api.loadbalancer.server.port=8080
    restart: unless-stopped

  web:
    image: ghcr.io/OWNER/urja-web:latest
    networks: [proxy, internal]
    environment:
      - NEXT_PUBLIC_API_URL=https://api.${DOMAIN}
    labels:
      - traefik.enable=true
      - traefik.http.routers.urja-web.rule=Host(`${DOMAIN}`)
      - traefik.http.routers.urja-web.tls.certresolver=letsencrypt
      - traefik.http.services.urja-web.loadbalancer.server.port=3000
    restart: unless-stopped

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: ${DB_USER}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME}
    volumes: [postgres_data:/var/lib/postgresql/data]
    networks: [internal]
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER}"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    volumes: [redis_data:/data]
    networks: [internal]
    healthcheck:
      test: ["CMD", "redis-cli", "ping"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

networks:
  proxy:
    external: true
  internal:
    driver: bridge

volumes:
  postgres_data:
  redis_data:
```

## Traefik Reference Setup (shared, not in Urja repo)

For VPS setup. One instance serves all apps.

**traefik/docker-compose.yml:**
- Traefik v3 image
- Ports 80, 443 exposed
- Creates the `proxy` network
- Mounts Docker socket (read-only) for auto-discovery
- acme.json volume for Let's Encrypt certs

**traefik/traefik.yml:**
- Entrypoints: web (:80), websecure (:443)
- HTTP → HTTPS redirect
- Let's Encrypt TLS resolver (tlsChallenge)
- Docker provider enabled, `exposedByDefault: false`
- Dashboard disabled

## CI/CD: GitHub Actions

**Trigger:** Push to `main`

**`.github/workflows/deploy.yml`:**

```
Job 1: build (parallel matrix)
  - Build API image → push ghcr.io/OWNER/urja-api:sha,latest
  - Build Web image → push ghcr.io/OWNER/urja-web:sha,latest

Job 2: deploy (after build)
  - SSH into VPS
  - cd /opt/urja
  - docker compose -f docker-compose.prod.yml pull
  - docker compose -f docker-compose.prod.yml up -d
  - Health check: curl https://api.domain.com/healthz
```

**GitHub Secrets needed:**
- `VPS_HOST` — VPS IP or hostname
- `VPS_SSH_KEY` — SSH private key for deploy user
- `VPS_USER` — SSH username (e.g., `deploy`)
- `DOMAIN` — Production domain

**Image tagging:** Each image tagged with both `latest` and git SHA. Rollback = pin to previous SHA.

## Files to Create/Modify

| File | Action | Purpose |
|------|--------|---------|
| `Dockerfile` | Modify | Add migrate binary + entrypoint script |
| `scripts/entrypoint.sh` | Create | Run migrations then exec API |
| `web/Dockerfile` | Create | Next.js standalone production image |
| `web/next.config.mjs` | Modify | Add `output: "standalone"` |
| `docker-compose.prod.yml` | Create | Production compose with Traefik labels |
| `.dockerignore` | Create | Exclude unnecessary files from build context |
| `web/.dockerignore` | Create | Exclude node_modules, .next from web context |
| `.github/workflows/deploy.yml` | Create | CI/CD pipeline |
| `docs/deploy/traefik-reference.md` | Create | Traefik VPS setup guide |

## Out of Scope

- Flutter mobile (native app, not containerized)
- Database backups (deferred)
- Monitoring/logging stack
- Multi-replica / load balancing (single VPS)
