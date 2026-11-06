# Database network access

Security maintenance by Codex, session /root, 2026-11-06 UTC on cdjk.
Redis and PostgreSQL host ports bind to loopback. Containers retain their internal network access.
The old Redis container referenced a removed polecat checkout. `docker-compose.redis.yml`
provides a standalone lifecycle using the existing `urja_redis_data` volume and `urja_default`
network: `PYTHONNOUSERSITE=1 docker-compose -p urja -f docker-compose.redis.yml up -d redis`.
No volume or application data is deleted. The main compose file needs the project's local `.env`.

The environment flag avoids incompatible user-installed Python dependencies in legacy Compose.
