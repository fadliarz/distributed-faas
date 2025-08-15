# Use a Debian-based Go image which uses glibc
FROM golang:1.24.4-bookworm AS builder

# Install build tools. Debian uses apt-get.
# libc6-dev provides the C standard library development files.
RUN apt-get update && apt-get install -y gcc libc6-dev

# Create the directory structure to match the replace directive (../../common)
WORKDIR /app/services/charge-service

# Copy go mod files
COPY services/charge-service/go.mod services/charge-service/go.sum ./

# Copy the common module (needed for local replacement) to the correct relative path
COPY common/ ../../common/

# Copy the infrastructure module (needed for local replacement) to the correct relative path
COPY infrastructure/kafka ../../infrastructure/kafka/

# Download dependencies
RUN go mod download

# Copy source code
COPY services/charge-service/ .

# Build the application
# The pre-compiled librdkafka library expects a glibc environment, which this image provides.
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
# Use a slim Debian image for a smaller final container
FROM debian:bookworm-slim

# Install ca-certificates for making HTTPS requests
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/services/charge-service/main .

# Expose port
EXPOSE 50059

# Run the binary
CMD ["./main"]