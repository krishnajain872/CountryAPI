# ===========================
# Stage 1: Build the Go binary
# ===========================
FROM golang:1.25-alpine AS builder

# Install necessary packages
RUN apk add --no-cache git bash

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first (for caching dependencies)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the entire project
COPY . .

# Build the Go binary
RUN go build -o /app/bin/server ./cmd/server/

# ===========================
# Stage 2: Create minimal runtime image
# ===========================
FROM alpine:3.20

# Install ca-certificates for HTTPS requests
RUN apk add --no-cache ca-certificates

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/bin/server .

# Copy any other required files (optional)
COPY --from=builder /app/.env .env

# Expose port
EXPOSE 3000

# Run the server
ENTRYPOINT ["./server"]
