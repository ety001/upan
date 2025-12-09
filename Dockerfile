# Build backend
FROM golang:1.23-alpine AS backend-builder

# Install build dependencies for CGO
RUN apk add --no-cache gcc musl-dev

WORKDIR /app/backend

# Copy go mod files
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy source code
COPY backend/ ./

# Update go.mod and build the application
RUN mkdir -p bin && go mod tidy && CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o bin/server ./cmd/server

# Build frontend
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

# Copy package files
COPY frontend/package.json frontend/package-lock.json* ./

# Install dependencies
RUN npm ci

# Copy source code
COPY frontend/ ./

# Build the application
RUN npm run build

# Final stage
FROM alpine:3.19

# Install dependencies
RUN apk --no-cache add \
    ca-certificates \
    sqlite \
    nginx \
    supervisor \
    curl

WORKDIR /app

# Copy backend binary
COPY --from=backend-builder /app/backend/bin/server ./backend/bin/server

# Copy frontend build
COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

# Copy nginx configuration
COPY docker-config/nginx/app.conf /etc/nginx/http.d/default.conf

# Copy supervisor configuration
COPY docker-config/supervisor/supervisord.conf /etc/supervisord.conf

# Create storage directories
RUN mkdir -p /app/storage/files /app/storage/db /var/log/supervisor

# Create nginx directories
RUN mkdir -p /var/log/nginx /var/cache/nginx /var/run/nginx

# Set permissions
RUN chmod +x /app/backend/bin/server

# Expose ports
EXPOSE 80

# Start supervisor
CMD ["/usr/bin/supervisord", "-c", "/etc/supervisord.conf"]

