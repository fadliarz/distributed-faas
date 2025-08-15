.PHONY: all compile clean generate-mocks

#
# PROTO GEN
#

gen-proto:
	make gen-tests-integration
	make gen-function-service
	make gen-invocation-service
	make gen-dispatcher-service
	make gen-machine
	make gen-registrar-service
	make gen-user-service
	make gen-charge-service

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

	protoc --proto_path=./proto \
	--go_out=./tests/integration/gen/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./tests/integration/gen/go \
	--go-grpc_opt=paths=source_relative \
	./proto/user-service/v1/api.proto

	protoc --proto_path=./proto \
	--go_out=./tests/integration/gen/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./tests/integration/gen/go \
	--go-grpc_opt=paths=source_relative \
	./proto/charge-service/v1/api.proto

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

	mkdir -p ./services/invocation-service/gen/go

	protoc --proto_path=./proto \
	--go_out=./services/invocation-service/gen/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./services/invocation-service/gen/go \
	--go-grpc_opt=paths=source_relative \
	./proto/invocation-service/v1/api.proto

	@echo "Protobuf compilation complete."

gen-dispatcher-service:
	@echo "Compiling Protobuf definitions for dispatcher service..."

	mkdir -p ./services/dispatcher-service/gen/go

	protoc --proto_path=./proto \
	--go_out=./services/dispatcher-service/gen/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./services/dispatcher-service/gen/go \
	--go-grpc_opt=paths=source_relative \
	./proto/machine-service/v1/api.proto

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

gen-user-service:
	@echo "Compiling Protobuf definitions for user service..."

	mkdir -p ./services/user-service/gen/go

	protoc --proto_path=./proto \
	--go_out=./services/user-service/gen/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./services/user-service/gen/go \
	--go-grpc_opt=paths=source_relative \
	./proto/user-service/v1/api.proto

	@echo "Protobuf compilation complete."

gen-charge-service:
	@echo "Compiling Protobuf definitions for charge service..."

	mkdir -p ./services/charge-service/gen/go

	protoc --proto_path=./proto \
	--go_out=./services/charge-service/gen/go \
	--go_opt=paths=source_relative \
	--go-grpc_out=./services/charge-service/gen/go \
	--go-grpc_opt=paths=source_relative \
	./proto/charge-service/v1/api.proto

	@echo "Protobuf compilation complete."

#
# Test
# 

test-create-invocation:
	@echo "Running integration tests for CDC function"

	yes | sudo rm -rf /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/zookeeper

	yes | sudo rm -rf /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/kafka

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/zookeeper/data

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/zookeeper/transactions

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/kafka/broker-1

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/kafka/dlq-1

	sudo chown -R 1000:1000 /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes

	go test ./tests/integration/create-invocation

test-create-invocation-verbose:
	@echo "Running integration tests for CDC function"

	yes | sudo rm -rf /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/zookeeper

	yes | sudo rm -rf /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/kafka

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/zookeeper/data

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/zookeeper/transactions

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/kafka/broker-1

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/kafka/dlq-1

	sudo chown -R 1000:1000 /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes

	go test -v ./tests/integration/create-invocation

test-create-user-verbose:
	@echo "Running integration tests for Create User"

	@echo "Removing existing Docker images..."
	-docker rmi distributed-faas-distributed-faas-user-processor:latest 2>/dev/null || true
	-docker rmi distributed-faas-distributed-faas-user-service:latest 2>/dev/null || true

	yes | sudo rm -rf /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/zookeeper

	yes | sudo rm -rf /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/kafka

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/zookeeper/data

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/zookeeper/transactions

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/kafka/broker-1

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/kafka/dlq-1

	sudo chown -R 1000:1000 /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes

	go test -v ./tests/integration/create-user

test-create-charge-verbose:
	@echo "Running integration tests for Create Charge"

	@echo "Removing existing Docker images..."
	-docker rmi distributed-faas-distributed-faas-charge-service:latest 2>/dev/null || true
	-docker rmi distributed-faas-distributed-faas-accumulator-service:latest 2>/dev/null || true

	yes | sudo rm -rf /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/zookeeper

	yes | sudo rm -rf /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/kafka

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/zookeeper/data

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/zookeeper/transactions

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/kafka/broker-1

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/kafka/dlq-1

	sudo chown -R 1000:1000 /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes

	go test -v ./tests/integration/create-charge

test-calculate-billing-verbose:
	@echo "Running integration tests for Calculate Billing"

	@echo "Removing existing Docker images..."
	-docker rmi distributed-faas-distributed-faas-billing-calculator-service:latest 2>/dev/null || true

	yes | sudo rm -rf /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/zookeeper

	yes | sudo rm -rf /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/kafka

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/zookeeper/data

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/zookeeper/transactions

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/kafka/broker-1

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/kafka/dlq-1

	sudo chown -R 1000:1000 /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes

	go test -v ./tests/integration/calculate-billing