# Multi-stage build for WhatsApp Proxy Go
# Stage 1: Build
FROM golang:1.22-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build binary
ARG VERSION=dev
ARG BUILD_TIME
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME}" \
    -o whatsapp-proxy \
    ./cmd/whatsapp-proxy

# Stage 2: Runtime
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1000 whatsapp && \
    adduser -D -u 1000 -G whatsapp whatsapp

# Create directories
RUN mkdir -p /etc/whatsapp-proxy /data && \
    chown -R whatsapp:whatsapp /etc/whatsapp-proxy /data

# Copy binary from builder
COPY --from=builder /build/whatsapp-proxy /usr/local/bin/whatsapp-proxy
RUN chmod +x /usr/local/bin/whatsapp-proxy

# Copy example config
COPY configs/config.example.yaml /etc/whatsapp-proxy/config.example.yaml

# Switch to non-root user
USER whatsapp

# Expose ports
EXPOSE 8443 8199

# Volume for config and data
VOLUME ["/etc/whatsapp-proxy", "/data"]

# Health check
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
    CMD wget -q -O- http://localhost:8199/health || exit 1

# Run proxy
ENTRYPOINT ["/usr/local/bin/whatsapp-proxy"]
CMD ["--config", "/etc/whatsapp-proxy/config.yaml"]
