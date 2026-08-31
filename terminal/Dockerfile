FROM golang:1.26.7-alpine AS builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o portfolio ./cmd/portfolio

FROM alpine:3.23
RUN apk add --no-cache ca-certificates netcat-openbsd \
    && addgroup -g 10001 portfolio \
    && adduser -D -u 10001 -G portfolio portfolio \
    && mkdir -p /app/.ssh && chown portfolio:portfolio /app/.ssh
WORKDIR /app
COPY --from=builder /build/portfolio /app/portfolio
USER 10001:10001
ENV SSH_PORT=2222
EXPOSE 2222
HEALTHCHECK --interval=15s --timeout=3s --start-period=5s --retries=3 \
    CMD nc -z 127.0.0.1 "$SSH_PORT" || exit 1
ENTRYPOINT ["/app/portfolio"]
