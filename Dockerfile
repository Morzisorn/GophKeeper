# ==============================================================================
# Multi-stage Dockerfile for GophKeeper
# Builds both server and agent (client) applications
# ==============================================================================

# Build stage - compiles Go applications
FROM golang:1.23.5-alpine AS builder

WORKDIR /app

# Download dependencies first (better Docker layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and build both applications
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o gophkeeper cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -o agent cmd/agent/main.go

# ==============================================================================
# Server image (for deployment/gophkeeper in k8s)
# ==============================================================================
FROM alpine:latest AS server
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the server binary from builder stage
COPY --from=builder /app/gophkeeper .

# Copy RSA keys for encryption
RUN mkdir -p ./keys
COPY --from=builder /app/*.pem ./keys/

EXPOSE 8080
CMD ["./gophkeeper"]

# ==============================================================================
# Agent image (for agent pod in k8s)
# ==============================================================================
FROM alpine:latest AS agent
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy the agent binary from builder stage
COPY --from=builder /app/agent .

CMD ["./agent"]
