# Use Go 1.23 alpine image
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files first
COPY services/execution/go.mod services/execution/go.sum ./

# Copy shared logger to match replace directive path (../../shared/logger from /app = /shared/logger)
COPY shared/logger /shared/logger

# Download dependencies
RUN go mod download

# Copy source code
COPY services/execution .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/main .

# Expose port
EXPOSE 8003

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8003/health || exit 1

# Run the application
CMD ["./main"]

