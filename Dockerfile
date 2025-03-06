# Start with the official Golang image
FROM golang:1.23.1 AS builder

# Enable CGO explicitly
ENV CGO_ENABLED=1

# Set the working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Install necessary system dependencies for Kafka (librdkafka)
RUN apt-get update && apt-get install -y librdkafka-dev pkg-config gcc g++ && rm -rf /var/lib/apt/lists/*

# Copy the source code
COPY . .

# Ensure all dependencies are available
RUN go mod tidy
RUN go mod vendor  # If using vendor mode

# Build the Go application with CGO enabled
RUN GOOS=linux GOARCH=amd64 go build -o tesla_server . && ls -lah tesla_server

# Use a lightweight base image for production
FROM alpine:latest

# Install required runtime dependencies
RUN apk add --no-cache ca-certificates libc6-compat librdkafka

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder
COPY --from=builder /app/tesla_server .

# Expose the port your app runs on (e.g., 8099)
EXPOSE 8099

# Run the application
CMD ["./tesla_server"]
