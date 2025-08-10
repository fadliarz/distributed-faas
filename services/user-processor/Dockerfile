FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first for dependency caching
COPY go.work go.work.sum ./
COPY common/go.mod common/go.sum ./common/
COPY infrastructure/kafka/go.mod infrastructure/kafka/go.sum ./infrastructure/kafka/
COPY services/user-processor/go.mod services/user-processor/go.sum ./services/user-processor/

# Download dependencies
RUN go work sync
RUN go mod download

# Copy source code
COPY common/ ./common/
COPY infrastructure/kafka/ ./infrastructure/kafka/
COPY services/user-processor/ ./services/user-processor/

# Build the application
WORKDIR /app/services/user-processor
RUN go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

COPY --from=builder /app/services/user-processor/main .

CMD ["./main"]
