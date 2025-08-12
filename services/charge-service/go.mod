module github.com/fadliarz/distributed-faas/services/charge-service

go 1.24.4

require (
	github.com/confluentinc/confluent-kafka-go/v2 v2.11.0
	github.com/fadliarz/distributed-faas/common v0.0.0
	github.com/golang/protobuf v1.5.4
	github.com/joho/godotenv v1.5.1
	github.com/rs/zerolog v1.34.0
	google.golang.org/grpc v1.74.2
)

require (
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/stretchr/testify v1.10.0 // indirect
	go.opentelemetry.io/otel v1.37.0 // indirect
	golang.org/x/net v0.42.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace github.com/fadliarz/distributed-faas/common => ../../common

replace github.com/fadliarz/distributed-faas/infrastructure/kafka => ../../infrastructure/kafka
