FROM golang:1.24.4-alpine AS builder

# Create the directory structure to match the replace directive (../../common)
WORKDIR /app/services/retry-service

# Copy go mod files
COPY services/retry-service/go.mod services/retry-service/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY services/retry-service/ .

# Build the application ---1
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/services/retry-service/main .

# Run the binary
CMD ["./main"]