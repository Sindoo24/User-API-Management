# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Install sqlc
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.30.0

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code and SQL files
COPY . .

# Generate SQLC code
RUN sqlc generate

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/server ./cmd/server

# Runtime stage
FROM alpine:latest

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates

# Copy binary from builder
COPY --from=builder /app/server .

# Copy migrations
COPY --from=builder /app/db/migrations ./db/migrations

EXPOSE 8080

CMD ["./server"]
