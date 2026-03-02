# ── Stage 1: Build Vue.js frontend ──────────────────────────────────────────
FROM node:22-alpine AS frontend-builder

WORKDIR /app/frontend

# Install dependencies first (layer cache)
COPY frontend/package*.json ./
RUN npm ci

# Copy source and build
COPY frontend/ ./
RUN npm run build

# ── Stage 2: Build Go backend ─────────────────────────────────────────────────
FROM golang:1.23-alpine AS backend-builder

WORKDIR /app/backend

# Install dependencies first (layer cache)
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy backend source
COPY backend/ ./

# Copy built frontend into backend/dist for embedding
COPY --from=frontend-builder /app/frontend/dist ./dist

# Build the Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/personal-coach .

# ── Stage 3: Minimal runtime image ────────────────────────────────────────────
FROM alpine:3.21

# Add ca-certificates for HTTPS calls to Anthropic API
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Copy the binary
COPY --from=backend-builder /app/personal-coach .

# Non-root user for security
RUN addgroup -g 1001 coach && adduser -u 1001 -G coach -s /bin/sh -D coach
RUN chown -R coach:coach /app
USER coach

EXPOSE 8080

CMD ["./personal-coach"]
