# ---- 1. Build stage ----
FROM golang:1.26-alpine AS builder

# Install git (needed for private modules sometimes)
RUN apk add --no-cache git

WORKDIR /app

# Cache dependencies first
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary (cmd/api is your entrypoint)
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o server .

# ---- 2. Runtime stage ----
FROM alpine:3.20

WORKDIR /app

# Add certs (important for HTTPS calls, AWS, etc.)
RUN apk add --no-cache curl ca-certificates

# Copy binary from builder
COPY --from=builder /app/server .

# Expose port (adjust if needed)
EXPOSE 8000

# Run the binary
CMD ["./server"]