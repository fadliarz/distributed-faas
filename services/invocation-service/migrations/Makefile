create-migration:
	@echo "Creating new migration..."
	migrate create -ext json -dir ./ -seq $(name)
	@echo "Migration created successfully."

dev-migrate-up:
	@echo "Running migrations... $(name=)"
	migrate -path ./ -database "mongodb://admin:password@localhost:27017/invocation-service-db?authSource=admin&replicaSet=rs0&directConnection=true" up
	@echo "Migrations completed successfully."

dev-migrate-down:
	@echo "Running migrations... $(name=)"
	migrate -path ./ -database "mongodb://admin:password@localhost:27017/invocation-service-db?authSource=admin&replicaSet=rs0&directConnection=true" down
	@echo "Migrations completed successfully."  