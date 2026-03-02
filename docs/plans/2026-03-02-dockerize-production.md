# Dockerize Urja for Production — Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Production-ready Docker setup for Urja on a single VPS with shared Traefik proxy and GitHub Actions CI/CD.

**Architecture:** Shared Traefik reverse proxy on an external `proxy` Docker network serves all apps. Urja runs api + web containers (on `proxy` + `internal` networks) with postgres + redis on `internal` only. GitHub Actions builds images to ghcr.io on push to main, then SSH-deploys to the VPS.

**Tech Stack:** Docker, Docker Compose, Traefik v3, GitHub Actions, GitHub Container Registry (ghcr.io), golang-migrate, Next.js standalone output.

**Repo:** `github.com/krantiutils/urja` (origin remote)

---

### Task 1: Create `.dockerignore`

**Files:**
- Create: `.dockerignore`

**Step 1: Create the file**

```
# .dockerignore
.git
.github
.env
.env.*
*.png
*.md
web/
mobile/
flutter/
firmware/
docs/
tests/
bin/
node_modules/
.next/
.playwright-mcp/
state.json
agent.md
coverage*
```

This keeps the Go API build context small. The web container uses its own context (`web/`).

**Step 2: Commit**

```bash
git add .dockerignore
git commit -m "chore: add .dockerignore for API build context"
```

---

### Task 2: Create `web/.dockerignore`

**Files:**
- Create: `web/.dockerignore`

**Step 1: Create the file**

```
# web/.dockerignore
node_modules/
.next/
screenshots/
tests/
*.md
playwright.config.ts
playwright-report/
```

**Step 2: Commit**

```bash
git add web/.dockerignore
git commit -m "chore: add web/.dockerignore for Next.js build context"
```

---

### Task 3: Modify `web/next.config.mjs` — enable standalone output

**Files:**
- Modify: `web/next.config.mjs`

**Step 1: Update the config**

Change from:
```js
/** @type {import('next').NextConfig} */
const nextConfig = {};

export default nextConfig;
```

To:
```js
/** @type {import('next').NextConfig} */
const nextConfig = {
  output: "standalone",
};

export default nextConfig;
```

**Step 2: Verify it builds locally**

Run: `cd web && npm run build`
Expected: Build succeeds, `.next/standalone/` directory is created.

**Step 3: Commit**

```bash
git add web/next.config.mjs
git commit -m "feat: enable Next.js standalone output for Docker"
```

---

### Task 4: Create `web/Dockerfile`

**Files:**
- Create: `web/Dockerfile`

**Step 1: Create the Dockerfile**

```dockerfile
### Dependencies
FROM node:20-alpine AS deps
WORKDIR /app
COPY package.json package-lock.json ./
RUN npm ci --ignore-scripts

### Builder
FROM node:20-alpine AS builder
WORKDIR /app
COPY --from=deps /app/node_modules ./node_modules
COPY . .

ARG NEXT_PUBLIC_API_URL
ENV NEXT_PUBLIC_API_URL=${NEXT_PUBLIC_API_URL}

RUN npm run build

### Runner
FROM node:20-alpine AS runner
WORKDIR /app

ENV NODE_ENV=production
ENV NEXT_TELEMETRY_DISABLED=1

RUN addgroup --system --gid 1001 nodejs
RUN adduser --system --uid 1001 nextjs

COPY --from=builder /app/public ./public
COPY --from=builder --chown=nextjs:nodejs /app/.next/standalone ./
COPY --from=builder --chown=nextjs:nodejs /app/.next/static ./.next/static

USER nextjs

EXPOSE 3000

ENV PORT=3000
ENV HOSTNAME="0.0.0.0"

CMD ["node", "server.js"]
```

**Step 2: Test the build**

Run: `docker build -t urja-web-test --build-arg NEXT_PUBLIC_API_URL=http://localhost:8080 web/`
Expected: Build succeeds.

**Step 3: Quick smoke test**

Run: `docker run --rm -p 3000:3000 urja-web-test`
Expected: Responds on http://localhost:3000

**Step 4: Clean up test image**

Run: `docker rmi urja-web-test`

**Step 5: Commit**

```bash
git add web/Dockerfile
git commit -m "feat: add Next.js production Dockerfile"
```

---

### Task 5: Modify API `Dockerfile` — add auto-migration

**Files:**
- Modify: `Dockerfile`
- Create: `scripts/entrypoint.sh`

**Step 1: Create the entrypoint script**

```bash
#!/bin/sh
set -e

echo "Running database migrations..."
migrate -path /app/migrations -database "$DATABASE_URL" up

echo "Starting urja-api..."
exec ./urja-api
```

Note: `DATABASE_URL` must be a full postgres connection string like `postgres://user:pass@host:5432/dbname?sslmode=disable`. The API itself reads individual `DB_*` env vars, but migrate CLI needs the URL form. The compose file will set both.

**Step 2: Update the Dockerfile**

Replace the entire `Dockerfile` with:

```dockerfile
### Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /app/urja-api ./cmd/api

### Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl

# Install golang-migrate
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-amd64.tar.gz | tar xz \
    && mv migrate /usr/local/bin/migrate

RUN addgroup -S urja && adduser -S urja -G urja

WORKDIR /app

COPY --from=builder /app/urja-api .
COPY --from=builder /app/db/migrations ./migrations
COPY scripts/entrypoint.sh ./entrypoint.sh
RUN chmod +x ./entrypoint.sh

USER urja

EXPOSE 8080

ENTRYPOINT ["./entrypoint.sh"]
```

Key changes from existing Dockerfile:
- Added `curl` to runtime for downloading migrate
- Install `golang-migrate` v4.17.0 binary
- Copy migrations from `db/migrations` (was `migrations` — verify this is correct path)
- Copy and use entrypoint script instead of direct binary exec

**Step 3: Test the build**

Run: `docker build -t urja-api-test .`
Expected: Build succeeds.

**Step 4: Clean up test image**

Run: `docker rmi urja-api-test`

**Step 5: Commit**

```bash
git add Dockerfile scripts/entrypoint.sh
git commit -m "feat: add auto-migration to API Dockerfile"
```

---

### Task 6: Create `docker-compose.prod.yml`

**Files:**
- Create: `docker-compose.prod.yml`

The existing `docker-compose.yml` stays as-is for local dev. This is the production file.

**Step 1: Create the file**

```yaml
services:
  api:
    image: ghcr.io/krantiutils/urja-api:${IMAGE_TAG:-latest}
    env_file: .env.prod
    environment:
      DATABASE_URL: postgres://${DB_USER:-urja}:${DB_PASSWORD}@postgres:5432/${DB_NAME:-urja}?sslmode=disable
      DB_HOST: postgres
      DB_PORT: "5432"
      REDIS_HOST: redis
      REDIS_PORT: "6379"
    depends_on:
      postgres:
        condition: service_healthy
      redis:
        condition: service_healthy
    networks:
      - proxy
      - internal
    labels:
      - traefik.enable=true
      - traefik.http.routers.urja-api.rule=Host(`api.${DOMAIN}`)
      - traefik.http.routers.urja-api.entrypoints=websecure
      - traefik.http.routers.urja-api.tls.certresolver=letsencrypt
      - traefik.http.services.urja-api.loadbalancer.server.port=8080
      - traefik.docker.network=proxy
    restart: unless-stopped

  web:
    image: ghcr.io/krantiutils/urja-web:${IMAGE_TAG:-latest}
    networks:
      - proxy
    labels:
      - traefik.enable=true
      - traefik.http.routers.urja-web.rule=Host(`${DOMAIN}`)
      - traefik.http.routers.urja-web.entrypoints=websecure
      - traefik.http.routers.urja-web.tls.certresolver=letsencrypt
      - traefik.http.services.urja-web.loadbalancer.server.port=3000
      - traefik.docker.network=proxy
    restart: unless-stopped

  postgres:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: ${DB_USER:-urja}
      POSTGRES_PASSWORD: ${DB_PASSWORD}
      POSTGRES_DB: ${DB_NAME:-urja}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - internal
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U ${DB_USER:-urja}"]
      interval: 5s
      timeout: 5s
      retries: 5
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    networks:
      - internal
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

Key differences from `docker-compose.yml` (dev):
- Uses pre-built images from ghcr.io (no `build:` directives)
- No ports exposed (Traefik handles external routing)
- Postgres/Redis on `internal` network only
- API/Web on both `proxy` + `internal`
- Traefik labels for automatic service discovery
- `DATABASE_URL` env var for migrate CLI; `DB_HOST`/`REDIS_HOST` overridden to container names
- `IMAGE_TAG` variable for pinning to specific commits (defaults to `latest`)

**Step 2: Commit**

```bash
git add docker-compose.prod.yml
git commit -m "feat: add production docker-compose with Traefik labels"
```

---

### Task 7: Create `.env.prod.example`

**Files:**
- Create: `.env.prod.example`

Template for the production env file that lives on the VPS (not committed with secrets).

**Step 1: Create the file**

```bash
# Domain (used by Traefik labels via docker-compose.prod.yml)
DOMAIN=urja.example.com

# Image tag (git SHA or "latest")
IMAGE_TAG=latest

# Database
DB_USER=urja
DB_PASSWORD=CHANGE_ME_STRONG_PASSWORD
DB_NAME=urja

# Server
SERVER_HOST=0.0.0.0
SERVER_PORT=8080
SERVER_READ_TIMEOUT=15s
SERVER_WRITE_TIMEOUT=15s
SERVER_SHUTDOWN_TIMEOUT=10s

# Auth
JWT_SECRET=CHANGE_ME_RANDOM_64_CHARS
ACCESS_TOKEN_EXPIRY=15m
REFRESH_TOKEN_EXPIRY=168h
BCRYPT_COST=12
OTP_EXPIRY=5m
OTP_COOLDOWN=60s
OTP_MAX_ATTEMPTS=5
OTP_HOURLY_LIMIT=10
DEV_OTP_BYPASS=false

# Redis
REDIS_PASSWORD=
REDIS_DB=0

# SMS (Aakash)
AAKASH_SMS_TOKEN=
AAKASH_SMS_API_URL=

# Khalti (payments)
KHALTI_SECRET_KEY=
KHALTI_PUBLIC_KEY=

# FCM (push notifications)
FCM_SERVICE_ACCOUNT_PATH=
```

**Step 2: Commit**

```bash
git add .env.prod.example
git commit -m "chore: add production env example file"
```

---

### Task 8: Create GitHub Actions deploy workflow

**Files:**
- Create: `.github/workflows/deploy.yml`

**Step 1: Create the workflow**

```yaml
name: Deploy

on:
  workflow_dispatch:
    inputs:
      image_tag:
        description: "Image tag to deploy (default: builds from HEAD)"
        required: false
        default: ""

env:
  REGISTRY: ghcr.io
  API_IMAGE: ghcr.io/krantiutils/urja-api
  WEB_IMAGE: ghcr.io/krantiutils/urja-web

jobs:
  build-api:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    steps:
      - uses: actions/checkout@v4

      - uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - uses: docker/build-push-action@v6
        with:
          context: .
          push: true
          tags: |
            ${{ env.API_IMAGE }}:latest
            ${{ env.API_IMAGE }}:${{ github.sha }}

  build-web:
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    steps:
      - uses: actions/checkout@v4

      - uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      - uses: docker/build-push-action@v6
        with:
          context: ./web
          push: true
          build-args: |
            NEXT_PUBLIC_API_URL=https://api.${{ secrets.DOMAIN }}
          tags: |
            ${{ env.WEB_IMAGE }}:latest
            ${{ env.WEB_IMAGE }}:${{ github.sha }}

  deploy:
    runs-on: ubuntu-latest
    needs: [build-api, build-web]
    steps:
      - name: Deploy to VPS
        uses: appleboy/ssh-action@v1
        with:
          host: ${{ secrets.VPS_HOST }}
          username: ${{ secrets.VPS_USER }}
          key: ${{ secrets.VPS_SSH_KEY }}
          script: |
            cd /opt/urja
            docker compose -f docker-compose.prod.yml pull
            docker compose -f docker-compose.prod.yml up -d --remove-orphans
            # Wait for API to be healthy
            sleep 5
            curl -sf http://localhost:8080/healthz || echo "WARNING: health check failed"
            # Clean up old images
            docker image prune -f
```

**GitHub Secrets required** (set in repo Settings > Secrets):
- `VPS_HOST` — VPS IP address
- `VPS_USER` — SSH username (e.g., `deploy`)
- `VPS_SSH_KEY` — SSH private key for that user
- `DOMAIN` — Production domain (e.g., `urja.example.com`)

**Step 2: Commit**

```bash
git add .github/workflows/deploy.yml
git commit -m "feat: add GitHub Actions deploy workflow for VPS"
```

---

### Task 9: Create Traefik reference setup doc

**Files:**
- Create: `docs/deploy/traefik-reference.md`

This is a reference doc for VPS setup — not part of the app itself, but needed for first-time server provisioning.

**Step 1: Create the doc**

````markdown
# Traefik Reverse Proxy Setup (VPS)

This is the shared Traefik instance that all apps connect to.
Set this up once on the VPS before deploying any app.

## Directory Structure

```
/opt/traefik/
├── docker-compose.yml
├── traefik.yml
└── acme.json          (auto-created, chmod 600)
```

## Setup

### 1. Create the proxy network

```bash
docker network create proxy
```

### 2. Create directory and files

```bash
sudo mkdir -p /opt/traefik
cd /opt/traefik
```

### 3. `docker-compose.yml`

```yaml
services:
  traefik:
    image: traefik:v3.3
    restart: unless-stopped
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock:ro
      - ./traefik.yml:/traefik.yml:ro
      - ./acme.json:/acme.json
    networks:
      - proxy

networks:
  proxy:
    external: true
```

### 4. `traefik.yml`

```yaml
api:
  dashboard: false

entryPoints:
  web:
    address: ":80"
    http:
      redirections:
        entryPoint:
          to: websecure
          scheme: https
  websecure:
    address: ":443"

certificatesResolvers:
  letsencrypt:
    acme:
      email: YOUR_EMAIL@example.com
      storage: /acme.json
      tlsChallenge: {}

providers:
  docker:
    endpoint: "unix:///var/run/docker.sock"
    exposedByDefault: false
    network: proxy
```

### 5. Create acme.json

```bash
touch acme.json
chmod 600 acme.json
```

### 6. Start Traefik

```bash
docker compose up -d
```

## Adding a New App

Each app just needs:
1. Join the `proxy` network (as `external: true`)
2. Add Traefik labels to web-facing services
3. `docker compose up -d`

Traefik auto-discovers services via Docker labels. No config reload needed.

## Deploying Urja

```bash
sudo mkdir -p /opt/urja
cd /opt/urja
# Copy docker-compose.prod.yml and .env.prod to this directory
docker compose -f docker-compose.prod.yml up -d
```
````

**Step 2: Commit**

```bash
git add docs/deploy/traefik-reference.md
git commit -m "docs: add Traefik VPS setup reference"
```

---

### Task 10: Verify full stack locally with docker compose

**Step 1: Build both images locally**

```bash
docker build -t urja-api-local .
docker build -t urja-web-local --build-arg NEXT_PUBLIC_API_URL=http://localhost:8080 web/
```

Expected: Both builds succeed.

**Step 2: Verify API image runs migrations**

```bash
docker compose up -d postgres redis
# Wait for healthy
sleep 5
docker run --rm --network urja_default \
  -e DATABASE_URL="postgres://urja:urja_secret@postgres:5432/urja?sslmode=disable" \
  -e DB_HOST=postgres -e DB_PORT=5432 -e DB_USER=urja -e DB_PASSWORD=urja_secret \
  -e DB_NAME=urja -e DB_SSLMODE=disable -e DB_MAX_CONNS=5 -e DB_MIN_CONNS=2 \
  -e REDIS_HOST=redis -e REDIS_PORT=6379 -e REDIS_DB=0 \
  -e JWT_SECRET=testsecret -e SERVER_HOST=0.0.0.0 -e SERVER_PORT=8080 \
  -p 8080:8080 \
  urja-api-local
```

Expected: See "Running database migrations..." then "server starting" in logs.

**Step 3: Verify health check**

```bash
curl http://localhost:8080/healthz
```

Expected: `{"status":"ok"}`

**Step 4: Clean up**

```bash
docker compose down
docker rmi urja-api-local urja-web-local
```

**Step 5: Final commit — all files together**

```bash
git push
```
