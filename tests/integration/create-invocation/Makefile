.PHONY: test

test:
	@echo "Running integration tests for CDC function"

	yes | sudo rm -rf /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/zookeeper

	yes | sudo rm -rf /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/kafka

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/zookeeper/data

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/zookeeper/transactions

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/kafka/broker-1

	mkdir -p /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes/kafka/dlq-1

	sudo chown -R 1000:1000 /home/fadlinux/workspace/distributed-faas/infrastructure/docker-compose/composes/volumes

	go test .	