FROM golang:1.24.4-alpine AS builder

# Create the directory structure to match the replace directive (../../common)
WORKDIR /app/services/user-service

# Copy go mod files
COPY services/user-service/go.mod services/user-service/go.sum ./

# Copy the common module (needed for local replacement) to the correct relative path
COPY common/ ../../common/

# Download dependencies
RUN go mod download

# Copy source code
COPY services/user-service/ .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/services/user-service/main .

# Expose port
EXPOSE 50050

# Run the binary
CMD ["./main"]
