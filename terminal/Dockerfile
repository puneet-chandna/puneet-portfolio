# ── Build Stage ──────────────────────────────────────────────
FROM golang:1.24-alpine AS builder

WORKDIR /build

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o portfolio ./cmd/portfolio

# ── Runtime Stage ────────────────────────────────────────────
FROM alpine:3.21

RUN apk add --no-cache ca-certificates netcat-openbsd

WORKDIR /app

# Copy binary and static assets
COPY --from=builder /build/portfolio .
COPY --from=builder /build/assets/ ./assets/

EXPOSE 22

ENTRYPOINT ["./portfolio"]
