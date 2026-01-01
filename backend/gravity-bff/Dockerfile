# ---- Build Stage ----
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Copy module files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application into a static binary
# CGO_ENABLED=0 for static linking
# -ldflags for smaller binary size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags="-w -s" \
    -a -installsuffix cgo \
    -o /app/gravity-bff \
    ./cmd/api/main.go

# ---- Production Stage ----
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS calls and tzdata for timezone support
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user for security
RUN adduser -D -g '' appuser

# Copy the compiled binary from the builder stage
COPY --from=builder /app/gravity-bff .

# Copy migrations for the migrate container
COPY --from=builder /app/migrations ./migrations

# Use non-root user
USER appuser

# Expose the application port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Set the entrypoint
CMD ["./gravity-bff"]
