FROM golang:1.24.4-alpine AS builder

# Create the directory structure to match the replace directive (../../common)
WORKDIR /app/services/registrar-service

# Copy go mod files
COPY services/registrar-service/go.mod services/registrar-service/go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY services/registrar-service/ .

# Build the application ---1
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/services/registrar-service/main .

# Expose port
EXPOSE 50050

# Run the binary
CMD ["./main"]