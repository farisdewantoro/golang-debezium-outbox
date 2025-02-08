# Use the official Golang image
FROM golang:1.22-alpine AS builder

# Install necessary build dependencies for CGO (glibc)
RUN apk update && apk add --no-cache \
    build-base \
    bash \
    libgcc \
    libc-dev

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum first
COPY go.mod go.sum ./

# Download dependencies
RUN go mod tidy

# Copy the rest of the source code
COPY . .

# Enable CGO and build the Go application with musl tag
RUN CGO_ENABLED=1 go build -tags musl -o /app/eventdriven-system .

# Use a minimal base image for production
FROM alpine:latest

# Install libc-dev for CGO support
RUN apk update && apk add --no-cache libc-dev

# Set working directory
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/eventdriven-system .

# Explicitly copy config.yaml (if available)
COPY config.yaml /app/config.yaml

# Ensure correct permissions
RUN chmod +x /app/eventdriven-system

# Set default entrypoint
ENTRYPOINT ["./eventdriven-system"]
