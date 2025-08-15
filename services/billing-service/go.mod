module github.com/fadliarz/distributed-faas/services/billing-service

go 1.24.4

require (
	github.com/fadliarz/distributed-faas/common v0.0.0-00010101000000-000000000000
	github.com/golang/protobuf v1.5.4
	github.com/joho/godotenv v1.5.1
	github.com/rs/zerolog v1.33.0
	go.mongodb.org/mongo-driver v1.17.4
	google.golang.org/grpc v1.69.4
)

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.19 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	golang.org/x/crypto v0.40.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sync v0.16.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.27.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20241015192408-796eee8c2d53 // indirect
	google.golang.org/protobuf v1.36.2 // indirect
)

replace github.com/fadliarz/distributed-faas/common => ../../common
