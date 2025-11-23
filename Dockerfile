# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with optimization flags
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o proxiff ./cmd/proxiff && \
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o sample-server ./example/servers

# Final stage
FROM alpine:3.20

# Install ca-certificates and wget for HTTPS support and health checks
RUN apk --no-cache add ca-certificates=20250911-r0 wget=1.24.5-r0 && \
    addgroup -g 1000 proxiff && \
    adduser -D -u 1000 -G proxiff proxiff

WORKDIR /app

# Copy binaries from builder
COPY --from=builder /build/proxiff .
COPY --from=builder /build/sample-server .

# Change ownership of the application files
RUN chown -R proxiff:proxiff /app

# Switch to non-root user
USER proxiff

# Expose default port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Default command (can be overridden)
CMD ["./proxiff"]
