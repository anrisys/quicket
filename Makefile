# ======================================== START OF MAKE FILE SET UP ========================================
# Load environment variables
-include ./docker/.env

# The main directory for Docker Compose files
COMPOSE_DIR := docker

# Services directories
USER_SERVICE_DIR := user-service

.PHONY: all up down check-env clean rebuild-all \
	build-base-image build-base-image-if-not-exists \
	build-user-service run-user-service migrate-user-up migrate-user-down migrate-user-version \

# Safety check for required variables
check-env:
	@if [ -z "${USER_MYSQL_PASSWORD}" ]; then \
		echo "Error: USER_MYSQL_PASSWORD is not set"; \
		exit 1; \
	fi
	@if [ -z "${USER_MYSQL_USER}" ]; then \
		echo "Error: USER_MYSQL_USER is not set"; \
		exit 1; \
	fi
	@echo "✅ Environment variables check passed"

# ======================================== END OF MAKE FILE SET UP ========================================

# ======================================== START OF GENERAL COMMANDS ========================================
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

# Clean up Docker images and volumes
clean: down
	@echo "Cleaning up Docker resources..."
	docker-compose -f $(COMPOSE_DIR)/docker-compose.dev.yml rm -fsv
	docker volume rm $(shell docker volume ls -q --filter name=quicket) || true

# Force rebuild everything
rebuild-all:
	docker system prune -f
	make build-base-image
	make build-user-service
# ======================================== END OF GENERAL COMMANDS ========================================

# ======================================== BASE IMAGE ========================================
# Build base image with better context and network handling
build-base-image: check-env
	@echo "Building the base development image..."
	docker build \
		--network=host \
		-t quicket-base-dev:latest \
		-f $(COMPOSE_DIR)/base-dev.Dockerfile 

# Alternative: Skip if image exists
build-base-image-if-not-exists: check-env
	@if docker image inspect quicket-base-dev:latest >/dev/null 2>&1; then \
		echo "Base image already exists, skipping build..."; \
	else \
		make build-base-image; \
	fi
# ======================================== END OF BASE IMAGE ========================================

# ======================================== START OF USER SERVICES ======================================== 
# Default values (optional safety net)
USER_MYSQL_USER ?= quicket_user
USER_MYSQL_PASSWORD ?= 
USER_MYSQL_DB_NAME ?= quicket_user_service

# Build user service using existing base image
build-user-service: build-base-image-if-not-exists
	@echo "Building user-service..."
	docker-compose -f $(COMPOSE_DIR)/docker-compose.dev.yml build user-api

# Start only the user-api and its dependencies
run-user-service: build-user-service
	@echo "Running user-service..."
	docker-compose -f $(COMPOSE_DIR)/docker-compose.dev.yml up -d user-api

# Migrate user services
migrate-user-up: check-env
	@echo "Running database migrations..."
	@echo "Database: ${USER_MYSQL_DB_NAME}, User: ${USER_MYSQL_USER}"
	docker run --rm \
		--network quicket-net \
		--env-file ./docker/.env \
		-v "./user-service/migrations:/migrations" \
		migrate/migrate \
		-path=/migrations \
		-database="mysql://${USER_MYSQL_USER}:${USER_MYSQL_PASSWORD}@tcp(${USER_MYSQL_HOST}:3306)/${USER_MYSQL_DB_NAME}" \
		up
	@echo "✅ Users Migrations completed successfully"

# Rollback user services
migrate-user-down: check-env
	@echo "Rolling back last migration..."
	docker run --rm \
		-v "./user-service/migrations:/migrations" \
		--network quicket-net \
		--env-file .env \
		migrate/migrate \
		-path=/migrations \
		-database="mysql://${USER_MYSQL_USER}:${USER_MYSQL_PASSWORD}@localhost:${USER_MYSQL_EXTERNAL_PORT}/${USER_MYSQL_DB_NAME}" \
		down
	@echo "✅ Users Migration rolled back"

migrate-user-version: check-env
	@echo "Current migration version:"
	docker run --rm \
		--network quicket-net \
		--env-file .env \
		-v "./user-service/migrations:/migrations" \
		migrate/migrate \
		-path=/migrations \
		-database="mysql://${USER_MYSQL_USER}:${USER_MYSQL_PASSWORD}@localhost:${USER_MYSQL_EXTERNAL_PORT}/${USER_MYSQL_DB_NAME}" \
		version

# ======================================== END OF USER SERVICES ======================================== 