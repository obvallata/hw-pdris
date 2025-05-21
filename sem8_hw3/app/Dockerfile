FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Use a smaller image for the final container
FROM alpine:latest

WORKDIR /app

# Create directory for logs
RUN mkdir -p /app/logs

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose the application port
EXPOSE 8000

# Run the application
CMD ["./main"] 