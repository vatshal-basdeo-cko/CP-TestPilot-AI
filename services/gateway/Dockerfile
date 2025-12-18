FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

# Copy go mod files first
COPY services/gateway/go.mod services/gateway/go.sum ./

# Copy shared logger to match replace directive path (../../shared/logger from /app = /shared/logger)
COPY shared/logger /shared/logger

# Download dependencies
RUN go mod download

# Copy source code
COPY services/gateway .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

FROM alpine:latest
RUN apk --no-cache add ca-certificates

WORKDIR /root/
COPY --from=builder /app/main .

EXPOSE 8000

HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8000/health || exit 1

CMD ["./main"]




