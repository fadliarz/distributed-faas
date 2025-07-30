start:
	./start-up.sh

shutdown:
	./shutdown.sh

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

mongo-up:
	@echo "Running docker-compose up for mongo..."
	docker-compose -f common.yml -f mongo.yml up -d
	@echo "docker-compose up complete."

mongo-down:
	@echo "Running docker-compose down for mongo..."
	docker-compose -f common.yml -f mongo.yml down
	@echo "docker-compose down complete."

inspect-mongo: 
	@echo "Inspecting MongoDB logs..."
	docker logs -f  docker-compose-distributed-faas-mongo-1 | grep -E -i "warn|warning|error|codeName"
	@echo "MongoDB logs inspection complete."


inspect-connect: 
	@echo "Inspecting Connect logs..."
	docker logs -f  docker-compose-distributed-faas-kafka-debezium-connect-1 | grep -E -i "warn|warning|error|codeName"
	@echo "Connect logs inspection complete."

show-topics: 
	@echo "Listing Kafka topics..."
	docker exec distributed-faas-kafka-broker-1 kafka-topics --bootstrap-server localhost:9092 --list 
	@echo "Kafka topics listed."

count-messages: 
	@echo "Counting messages in Kafka topic 'cdc.function-db.function'..."
	docker exec docker-compose-distributed-faas-kafka-broker-1-1 kafka-run-class kafka.tools.GetOffsetShell --broker-list localhost:9092 --topic cdc.function-db.function
	@echo "Message count in Kafka topic 'cdc.function-db.function' complete."

bash-mongo:	
	@echo "Opening bash shell in MongoDB container..."
	docker exec -it docker-compose-distributed-faas-mongo-1 /bin/bash
	@echo "Bash shell in MongoDB container opened."

build-nodejs:
	@echo "Building Node.js Docker image..."
	docker build -t nodejs-app -f Dockerfile.nodejs .
	@echo "Node.js Docker image built."