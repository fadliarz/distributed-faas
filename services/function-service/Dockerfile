FROM golang:1.24.4-alpine AS builder

# Create the directory structure to match the replace directive (../../common)
WORKDIR /app/services/function-service

# Copy go mod files
COPY services/function-service/go.mod services/function-service/go.sum ./

# Copy the common module (needed for local replacement) to the correct relative path
COPY common/ ../../common/

# Copy source code
COPY services/function-service/ .

# Tidy up go.mod and go.sum files
RUN go mod tidy

# Build the application ---1
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/services/function-service/main .

# Run the binary
CMD ["./main"]