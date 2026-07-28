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

# Which domain's subdomains are tenant gym sites. Next inlines NEXT_PUBLIC_*
# at build time — including into the middleware that does the routing — so this
# has to be a build arg. Setting it at runtime silently does nothing.
ARG NEXT_PUBLIC_BASE_DOMAIN=nepalgym.xyz
ENV NEXT_PUBLIC_BASE_DOMAIN=${NEXT_PUBLIC_BASE_DOMAIN}

# Ensure public directory exists (may be empty)
RUN mkdir -p public

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
