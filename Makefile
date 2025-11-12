# ========================================
# CONFIGURATION LOADING
# ========================================

-include ./docker/.env

# Directory configuration
COMPOSE_DIR := docker
USER_SERVICE_DIR := user-service

# ========================================
# PHONY TARGET DECLARATIONS
# ========================================

.PHONY: help check-env build-base-image build-base-image-if-not-exists \
        gateway-up gateway-logs gateway-restart \
        build-user-service run-user-service \
        migrate-user-up migrate-user-down migrate-user-version \
		personal-dev team-dev build-images down clean \
		build-booking-image build-api-gateway \
		restart-user-api logs-user-api user-shell-db \
		restart-rabbitmq rabbitmq-logs \
		restart-booking-api logs-booking-api booking-shell-db \

# ========================================
# HELP COMMAND
# ========================================

## help: Show this help message
help:
	@echo "QuickET Development Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Targets:"
	@awk '/^## / { \
		help_line = substr($$0, 4); \
		getline comment; \
		while (comment ~ /^# /) { \
			help_line = help_line "\n           " substr(comment, 3); \
			getline comment; \
		} \
		printf "  \033[36m%-20s\033[0m %s\n", $$2, help_line; \
	}' $(MAKEFILE_LIST) | sort

# ========================================
# ENVIRONMENT & VALIDATION
# ========================================

## check-env: Validate required environment variables
check-env:
	@if [ -z "${USER_MYSQL_PASSWORD}" ]; then \
		echo "Error: USER_MYSQL_PASSWORD is not set"; \
		exit 1; \
	fi
	@if [ -z "${USER_MYSQL_USER}" ]; then \
		echo "Error: USER_MYSQL_USER is not set"; \
		exit 1; \
	fi
	@if [ -z "${USER_MYSQL_DB_NAME}" ]; then \
		echo "Error: USER_MYSQL_DB_NAME is not set"; \
		exit 1; \
	fi
	@echo "✅ Environment variables check passed"

# ========================================
# BASE IMAGE MANAGEMENT
# ========================================

## build-base-image: Build the base development image
build-base-image:
	@echo "Building base image..."
	docker build \
		--network=host \
		-t quicket-base-dev:latest . \
		-f $(COMPOSE_DIR)/base-dev.Dockerfile 

# build-base-image-if-not-exists: Check build image existence before building it
build-base-image-if-not-exists: check-env
	@if docker image inspect quicket-base-dev:latest >/dev/null 2>&1; then \
		echo "Base image already exists, skipping build..."; \
	else \
		make build-base-image; \
	fi

# ========================================
# GATEWAY MANAGEMENT
# ========================================

## gateway-up: Start the gateway service
gateway-up: 
	@echo "Starting gateway..."
	cd docker && docker compose up -d api-gateway

## gateway-logs: Show gateway service logs
gateway-logs: 
	cd docker && docker compose logs -f api-gateway

## gateway-restart: Restart the gateway service
gateway-restart: 
	@echo "Restarting gateway..."
	docker compose -f docker/docker-compose.yml restart api-gateway

# ========================================
# MIGRATION COMMANDS
# ========================================

## migrate-user-up: Run user service database migrations up
migrate-user-up: check-env
	@echo "Running user migrations up..."
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

## migrate-user-down: Rollback user service database migrations
migrate-user-down: check-env
	@echo "Running user migrations down..."
	docker run --rm \
		-v "./user-service/migrations:/migrations" \
		--network quicket-net \
		--env-file .env \
		migrate/migrate \
		-path=/migrations \
		-database="mysql://${USER_MYSQL_USER}:${USER_MYSQL_PASSWORD}@localhost:${USER_MYSQL_EXTERNAL_PORT}/${USER_MYSQL_DB_NAME}" \
		down
	@echo "✅ Users Migration rolled back"

## migrate-user-version: Check current migration version
migrate-user-version: check-env
	@echo "Checking migration version..."
	docker run --rm \
		--network quicket-net \
		--env-file .env \
		-v "./user-service/migrations:/migrations" \
		migrate/migrate \
		-path=/migrations \
		-database="mysql://${USER_MYSQL_USER}:${USER_MYSQL_PASSWORD}@localhost:${USER_MYSQL_EXTERNAL_PORT}/${USER_MYSQL_DB_NAME}" \
		version

# ========================================
# USER-SERVICE-SPECIFIC COMMANDS
# ========================================

## build-user-service: Build the user service
build-user-service: build-base-image-if-not-exists
	@echo "Building user-service..."
	docker compose -f docker/docker-compose.override.yml build user-api

## run-user-service: Run the user service
run-user-service: build-user-service
	@echo "Running user-service..."
	docker compose -f docker/docker-compose.override.yml up -d user-api

# ========================================
# DEVELOPMENT TARGET
# ========================================

## run-dev: Fast startup with already created images separately
run-dev: check-env
	@echo "🚀 Running docker compose (make sure already have required images)..."
	cd docker && docker compose up -d

## personal-dev: Fast startup for personal development (uses pre-built)
build-and-run-dev: build-images
	@echo "🚀 Starting with pre-built images..."
	cd docker && docker compose up -d

## team-dev: Complete setup for team members (builds everything)
team-dev:
	@echo "🏗️ Building and starting for team development..."
	cd docker && docker compose -f docker/docker-compose.team.yml up --build -d

## build-images: Build all service images
build-images: 
	@echo "📦 Building service images..."
	# API Gateway (with a Dockerfile in /api-gateway directory)
	docker build -t ${API_GATEWAY_IMAGE_NAMETAG} ./api-gateway
	
	# Base development image
	docker build -t ${BASE_IMAGE_NAMETAG} -f docker/base-dev.Dockerfile .
	
	# User service
	docker build -t ${USER_IMAGE_NAMETAG} -f user-service/Dockerfile.dev ./user-service

	# Booking service
	docker build -t ${BOOKING_IMAGE_NAMETAG} -f booking-service/Dockerfile.dev ./booking-service

## build-booking-image: Build booking service image
build-booking-image: 
	@echo "📦 Building booking service image..."
	# Booking service
	docker build -t ${BOOKING_IMAGE_NAMETAG} -f booking-service/Dockerfile.dev ./booking-service

## build-api-gateway: Build api gateway image
build-api-gateway: 
	@echo "📦 Building api gateway image..."
	# API Gateway (with a Dockerfile in /api-gateway directory)
	docker build -t ${API_GATEWAY_IMAGE_NAMETAG} ./api-gateway

## down: Stop all services
down:
	cd docker && docker compose down

## clean: Stop and clean all services dan its data
clean:
	cd docker && docker compose down -v

# ========================================
# SERVICES MANAGEMENT
# ========================================

## restart-user-api: Restart the user API service
restart-user-api: 
	@echo "Restarting user API..."
	docker compose -f docker/docker-compose.yml restart user-api

## logs-user-api: Logs the user API service
user-api-logs: 
	cd docker && docker compose logs -f user-api

## user-shell-db: Go into bash shell of user-mysql
user-shell-db: 
	cd docker && docker compose exec -it user-mysql bash

## restart-booking-api: Restart the booking API service
restart-booking-api: 
	@echo "Restarting booking API..."
	docker compose -f docker/docker-compose.yml restart booking-api

## logs-booking-api: Logs the user API service
booking-api-logs: 
	cd docker && docker compose logs -f booking-api

## booking-shell-db: Go into bash shell of booking-mysql
booking-shell-db: 
	cd docker && docker compose exec -it booking-mysql bash

# ========================================
# MESSAGE BROKER
# ========================================

## restart-rabbitmq: Restart message broker
restart-rabbitmq: 
	@echo "Restarting rabbitmq..."
	docker compose -f docker/docker-compose.yml restart rabbitmq

## rabbitmq-logs: Logs the user message broker
rabbitmq-logs: 
	cd docker && docker compose logs -f rabbitmq