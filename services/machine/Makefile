.PHONY: build run clean

build:
	@echo "Building machine service..."
	go build -o bin/machine-service .
	@echo "Build complete."

run: build
	@echo "Starting machine service..."
	./bin/machine-service

clean:
	@echo "Cleaning up..."
	rm -rf bin/
	@echo "Clean complete."
