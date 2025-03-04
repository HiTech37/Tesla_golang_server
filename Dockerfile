# Start with the official Golang image
FROM golang:1.23.4 AS builder

# Set the working directory
WORKDIR /tesla_server

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the Go application
RUN go build -o main .

# Use a lightweight base image for production
FROM alpine:latest

# Set the working directory
WORKDIR /root/

# Copy the compiled binary from the builder
COPY --from=builder /tesla_server/main .

# Expose the port your app runs on (e.g., 8080)
EXPOSE 8099

# Run the application
CMD ["./main"]
