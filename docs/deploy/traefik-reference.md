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
