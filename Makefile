

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_aggregator build_notification build_subscriber
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_subscriber: builds the subscriber binary as a linux executable
build_subscriber:
	@echo "Building subscriber binary..."
	cd ./subscriber-service && env GOOS=linux CGO_ENABLED=0 go build -gcflags "all=-N -l" -o subscriberService .
	@echo "Done!"

## build_aggregator: builds the aggregator binary as a linux executable
build_aggregator:
	@echo "Building aggregator binary..."
	cd ./aggregator-service && env GOOS=linux CGO_ENABLED=0 go build -o aggregatorService .
	@echo "Done!"

## build_notification: builds the notification binary as a linux executable
build_notification:
	@echo "Building notification binary..."
	cd ./notification-service && env GOOS=linux CGO_ENABLED=0 go build -o notificationService .
	@echo "Done!"