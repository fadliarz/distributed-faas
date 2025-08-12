FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Copy common and infrastructure dependencies
COPY ../../common ./common/
COPY ../../infrastructure ./infrastructure/

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copy the binary
COPY --from=builder /app/main .

# Expose port
EXPOSE 50052

# Run the binary
CMD ["./main"]
