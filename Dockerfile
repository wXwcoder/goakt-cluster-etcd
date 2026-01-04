# Build stage
FROM golang:1.25.3-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/bin/accounts .

# Runtime stage
FROM alpine:latest

# Install ca-certificates for TLS connections
RUN apk --no-cache add ca-certificates

# Set working directory
WORKDIR /root/

# Copy binary from builder stage
COPY --from=builder /app/bin/accounts .

# Copy environment file if exists
COPY --from=builder /app/.env ./.env

# Expose the port the app runs on
EXPOSE 8080
EXPOSE 8555
EXPOSE 25520

# Command to run the application
CMD ["./accounts", "run"]