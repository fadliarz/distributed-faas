.PHONY: all compile clean generate-mocks

# Compile Protobuf definitions for services

gen-tests-integration:
	@echo "Compiling Protobuf definitions for integration tests..."

	mkdir -p ./tests/integration/gen/go

	protoc --proto_path=./proto \
	--go_out=./tests/integration/gen/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./tests/integration/gen/go \
	--go-grpc_opt=paths=source_relative \
	./proto/function-service/v1/api.proto

	protoc --proto_path=./proto \
	--go_out=./tests/integration/gen/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./tests/integration/gen/go \
	--go-grpc_opt=paths=source_relative \
	./proto/invocation-service/v1/api.proto

	protoc --proto_path=./proto \
	--go_out=./tests/integration/gen/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./tests/integration/gen/go \
	--go-grpc_opt=paths=source_relative \
	./proto/registrar-service/v1/api.proto

	@echo "Protobuf compilation complete."

gen-dispatcher-service:
	@echo "Compiling Protobuf definitions for dispatcher service..."
	protoc --proto_path=./proto \
	--go_out=./services/dispatcher-service/gen/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./services/dispatcher-service/gen/go \
	--go-grpc_opt=paths=source_relative \
	./proto/machine-service/v1/api.proto
	@echo "Protobuf compilation complete."

gen-function-service:
	@echo "Compiling Protobuf definitions for function service..."
	mkdir -p ./services/function-service/gen/go

	protoc --proto_path=./proto \
	--go_out=./services/function-service/gen/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./services/function-service/gen/go \
	--go-grpc_opt=paths=source_relative \
	./proto/function-service/v1/api.proto
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

gen-machine:
	@echo "Compiling Protobuf definitions for machine service..."

	mkdir -p ./services/machine/gen/go

	protoc --proto_path=./proto \
	--go_out=./services/machine/gen/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./services/machine/gen/go \
	--go-grpc_opt=paths=source_relative \
	./proto/machine-service/v1/api.proto

	@echo "Protobuf compilation complete."

gen-registrar-service:
	@echo "Compiling Protobuf definitions for registrar service..."

	mkdir -p ./services/registrar-service/gen/go

	protoc --proto_path=./proto \
	--go_out=./services/registrar-service/gen/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./services/registrar-service/gen/go \
	--go-grpc_opt=paths=source_relative \
	./proto/registrar-service/v1/api.proto

	@echo "Protobuf compilation complete."