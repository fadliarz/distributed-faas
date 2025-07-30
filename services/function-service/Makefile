# Migration Makefile for Function Service

.PHONY: migrate-up migrate-down migrate-status migrate-create migrate-build migrate-clean help

# Run unit tests
run-unit-test:
	@echo "Running unit tests..."
	@go test -v ./domain/application-service/service
	@echo "Unit tests completed"