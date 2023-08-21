# Makefile

# Build and start the Docker containers
start:
	docker-compose up -d

# Stop and remove the Docker containers
stop:
	docker-compose down

# Build the Docker images
build:
	docker-compose build

# View logs of the running containers
logs:
	docker-compose logs -f

# Execute a command inside the frontend container
frontend-shell:
	docker-compose exec frontend sh

# Execute a command inside the backend container
backend-shell:
	docker-compose exec backend sh

# Access the PostgreSQL database using psql
psql:
	docker-compose exec database psql -U postgres -d sensor_metadata

ui-test:
	cd sensor-metadata-ui && npm test

api-test:
	cd sensor-metadata-api && go test ./... -coverprofile=coverage.out

swagger:
	cd sensor-metadata-api && swag init --parseDependency

build-api:
	cd sensor-metadata-api && go build -o sensor-metadata-api

# Default target when no argument is provided
default: start
