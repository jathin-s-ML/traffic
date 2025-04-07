# Use the latest Go version as the builder stage
FROM golang:latest AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy Go modules files first for better caching
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application files
COPY . .

# Build the Go binary
RUN go build -o traffic-generator ./main.go

# Use a lightweight base image for the final container
FROM golang:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/traffic-generator /app/traffic-generator

# Copy the config file (it can also be mounted as a volume in docker-compose)
COPY config/config.yaml /app/config.yaml

# Set execution permissions (optional, but ensures the binary can run)
RUN chmod +x /app/traffic-generator

# Default command to run the application
CMD ["/app/traffic-generator"]
