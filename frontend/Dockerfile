# Build stage
FROM node:20-alpine AS builder

# Accept build-time environment variable for API Gateway URL
ARG VITE_API_GATEWAY_URL=http://localhost:9000
ENV VITE_API_GATEWAY_URL=${VITE_API_GATEWAY_URL}

WORKDIR /app

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm ci

# Copy source code
COPY . .

# Build the application (VITE_API_GATEWAY_URL is available at build time)
RUN npm run build

# Production stage
FROM nginx:alpine

# Copy built assets
COPY --from=builder /app/dist /usr/share/nginx/html

# Copy nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Expose port
EXPOSE 80

# Health check - use 127.0.0.1 to avoid IPv6 resolution issues
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://127.0.0.1:80/ || exit 1

CMD ["nginx", "-g", "daemon off;"]
