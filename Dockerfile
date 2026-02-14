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

RUN apk add --no-cache ca-certificates tzdata

RUN addgroup -S urja && adduser -S urja -G urja

WORKDIR /app

COPY --from=builder /app/urja-api .
COPY --from=builder /app/migrations ./migrations

USER urja

EXPOSE 8080

ENTRYPOINT ["./urja-api"]
