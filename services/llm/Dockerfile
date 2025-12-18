# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files first
COPY services/llm/go.mod services/llm/go.sum* ./

# Copy shared logger to match replace directive path (../../shared/logger from /app = /shared/logger)
COPY shared/logger /shared/logger

# Download dependencies
RUN go mod download

# Copy source code
COPY services/llm .

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o llm .

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates curl

# Copy binary from builder
COPY --from=builder /app/llm .

# Expose port
EXPOSE 8002

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8002/health || exit 1

# Run binary
CMD ["./llm"]

