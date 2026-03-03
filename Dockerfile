### Build stage
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/urja-api ./cmd/api

### Runtime stage
FROM alpine:3.19

RUN apk add --no-cache ca-certificates tzdata curl

# Install golang-migrate (arch-agnostic)
ARG TARGETARCH
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.17.0/migrate.linux-${TARGETARCH}.tar.gz | tar xz \
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
