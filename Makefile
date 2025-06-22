.PHONY: all compile clean

compile-function-service:
	@echo "Compiling Protobuf definitions for function service..."
	protoc --proto_path=./proto \
	--go_out=./services/function-service/gen/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./services/function-service/gen/go \
	--go-grpc_opt=paths=source_relative \
	./proto/function-service/v1/api.proto
	@echo "Protobuf compilation complete."
