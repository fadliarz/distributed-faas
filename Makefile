.PHONY: all compile clean generate-mocks

# Compile Protobuf definitions for services

gen-dispatcher-service:
	@echo "Compiling Protobuf definitions for dispatcher service..."
	protoc --proto_path=./proto \
	--go_out=./services/dispatcher-service/gen/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./services/dispatcher-service/gen/go \
	--go-grpc_opt=paths=source_relative \
	./proto/machine-service/v1/api.proto
	@echo "Protobuf compilation complete."

gen-machine:
	@echo "Compiling Protobuf definitions for machine..."
	protoc --proto_path=./proto \
	--go_out=./services/machine/gen/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./services/machine/gen/go \
	--go-grpc_opt=paths=source_relative \
	./proto/machine-service/v1/api.proto
	@echo "Protobuf compilation complete."

gen-invocation-service:
	@echo "Compiling Protobuf definitions for invocation service..."
	protoc --proto_path=./proto \
	--go_out=./services/invocation-service/gen/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./services/invocation-service/gen/go \
	--go-grpc_opt=paths=source_relative \
	./proto/invocation-service/v1/api.proto
	@echo "Protobuf compilation complete."
