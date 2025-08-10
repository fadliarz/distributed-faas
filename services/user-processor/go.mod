module github.com/fadliarz/distributed-faas/services/user-processor

go 1.24.4

require (
	github.com/fadliarz/distributed-faas/common v0.0.0-00010101000000-000000000000
	github.com/fadliarz/distributed-faas/infrastructure/kafka v0.0.0-00010101000000-000000000000
	github.com/joho/godotenv v1.5.1
	github.com/rs/zerolog v1.34.0
	go.mongodb.org/mongo-driver v1.17.4
)

require (
	cloud.google.com/go/compute/metadata v0.7.0 // indirect
	github.com/Azure/azure-sdk-for-go/sdk/internal v1.10.0 // indirect
	github.com/aws/aws-sdk-go-v2/feature/ec2/imds v1.16.32 // indirect
	github.com/aws/aws-sdk-go-v2/internal/ini v1.8.3 // indirect
	github.com/aws/aws-sdk-go-v2/service/sso v1.25.5 // indirect
	github.com/aws/aws-sdk-go-v2/service/ssooidc v1.30.3 // indirect
	github.com/aws/smithy-go v1.22.5 // indirect
	github.com/bufbuild/protocompile v0.8.0 // indirect
	github.com/compose-spec/compose-go/v2 v2.6.0 // indirect
	github.com/confluentinc/confluent-kafka-go/v2 v2.11.0 // indirect
	github.com/containerd/containerd/api v1.8.0 // indirect
	github.com/containerd/containerd/v2 v2.0.5 // indirect
	github.com/docker/cli v28.0.4+incompatible // indirect
	github.com/docker/compose/v2 v2.35.0 // indirect
	github.com/docker/docker v28.2.2+incompatible // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-viper/mapstructure/v2 v2.3.0 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/jhump/protoreflect v1.15.6 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/testcontainers/testcontainers-go v0.38.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/oauth2 v0.30.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto v0.0.0-20240325203815-454cdb8f5daa // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250528174236-200df99c418a // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace (
	github.com/fadliarz/distributed-faas/common => ../../common
	github.com/fadliarz/distributed-faas/infrastructure/kafka => ../../infrastructure/kafka
)
