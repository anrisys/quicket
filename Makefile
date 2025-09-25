# The main directory for Docker Compose files
COMPOSE_DIR := docker

# The directory for the user service
USER_SERVICE_DIR := user-service

.PHONY: all up down build-base-image \
	build-user-service run-user-service migrate-user-service clean

# Default target when no target is specified
all: up

# Build and start all services defined in docker-compose.dev.yml
up: build-base-image
	@echo "Starting all services..."
	docker-compose -f $(COMPOSE_DIR)/docker-compose.dev.yml up -d --build

# Stop and remove all services
down:
	@echo "Stopping all services..."
	docker-compose -f $(COMPOSE_DIR)/docker-compose.dev.yml down

# A new target to build and tag the reusable base image
build-base-image:
	@echo "Building the base development image..."
	docker build -t quicket-base-dev:latest -f $(COMPOSE_DIR)/base-dev.Dockerfile $(COMPOSE_DIR)

# Build the user-service-api image and its dependencies
build-user-service: build-base-image
	@echo "Building user-service..."
	docker-compose -f $(COMPOSE_DIR)/docker-compose.dev.yml build user-service-api

# Start only the user-service-api and its dependencies
run-user-service: build-user-service
	@echo "Running user-service..."
	docker-compose -f $(COMPOSE_DIR)/docker-compose.dev.yml up -d user-service-api

# Run the user-service database migrations
migrate-user-service:
	@echo "Running migrations for user-service..."
	docker-compose -f $(COMPOSE_DIR)/docker-compose.dev.yml up --build --force-recreate user-service-migrate

# Clean up Docker images and volumes
clean: down
	@echo "Cleaning up Docker resources..."
	docker-compose -f $(COMPOSE_DIR)/docker-compose.dev.yml rm -fsv
	docker volume rm $(shell docker volume ls -q --filter name=quicket) || true