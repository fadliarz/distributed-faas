# Script to manage Docker Compose

start:
	./start-up.sh

shutdown:
	./shutdown.sh

# Makefile for Docker Compose Management

up: 
	@echo "Running docker-compose up..."
	./start-up.sh
	@echo "docker-compose up complete."

down: 
	@echo "Running docker-compose down..."
	./shutdown.sh
	@echo "docker-compose down complete."

restart-services:
	@echo "Restarting services..."
	docker-compose -f ./composes/common.yml -f ./composes/services.yml --profile "*" down
	COMPOSE_BAKE=true docker-compose -f ./composes/common.yml -f ./composes/services.yml --profile "*" up -d --build
	@echo "Services restarted."

# Up and down related commands

mongo-up:
	@echo "Running docker-compose up for mongo..."
	docker-compose -f common.yml -f mongo.yml up -d
	@echo "docker-compose up complete."

mongo-down:
	@echo "Running docker-compose down for mongo..."
	docker-compose -f common.yml -f mongo.yml down
	@echo "docker-compose down complete."

# Inspect related commands

inspect-mongo: 
	@echo "Inspecting MongoDB logs..."
	docker logs -f  docker-compose-distributed-faas-mongo-1 | grep -E -i "warn|warning|error|codeName"
	@echo "MongoDB logs inspection complete."

inspect-connect: 
	@echo "Inspecting Connect logs..."
	docker logs -f  docker-compose-distributed-faas-kafka-debezium-connect-1 | grep -E -i "warn|warning|error|codeName"
	@echo "Connect logs inspection complete."

# Topics related commands

show-topics: 
	@echo "Listing Kafka topics..."
	docker exec distributed-faas-kafka-broker-1 kafka-topics --bootstrap-server localhost:9092 --list 
	@echo "Kafka topics listed."

count-messages: 
	@echo "Counting messages in Kafka topics..."
	
	docker exec distributed-faas-kafka-broker-1 kafka-run-class kafka.tools.GetOffsetShell --broker-list localhost:9092 --topic cdc.invocation-db.invocation
	
	docker exec distributed-faas-kafka-broker-1 kafka-run-class kafka.tools.GetOffsetShell --broker-list localhost:9092 --topic cdc.checkpoint-db.checkpoint
	
	@echo "Message count in Kafka topics complete."

test-create-user-count-messages:	
	@echo "Counting messages in Kafka topics..."

	docker exec distributed-faas-kafka-broker-1 kafka-run-class kafka.tools.GetOfsetShell --broker-list localhost:9092 --topic cdc.user-db.user

	@echo "Message count in Kafka topics complete."

# Execute related commands

bash-mongo:	
	@echo "Opening bash shell in MongoDB container..."
	docker exec -it docker-compose-distributed-faas-mongo mongosh -u admin -p password --authenticationDatabase admin
	@echo "Bash shell in MongoDB container opened."