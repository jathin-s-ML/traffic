# Use the latest Go version as the builder stage
FROM golang:latest AS builder

WORKDIR /app

# Copy Go modules files first for better caching
COPY go.mod go.sum ./
RUN go mod tidy

# Copy the rest of the application files
COPY . .

# Build the Go binary
RUN go build -o traffic-stats-col .

# Use a lightweight base image for the final container
FROM golang:latest

WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/traffic-stats-col /app/traffic-stats-col

# Copy the config file (it can also be mounted as a volume in docker-compose)
COPY config.yaml /app/config.yaml

# Set execution permissions (optional, but ensures the binary can run)
RUN chmod +x /app/traffic-stats-col

# Default command to run the application
CMD ["/app/traffic-stats-col"]
