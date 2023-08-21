# sensor-metadata-app

JSON REST API for storing and querying sensor metadata.

## Prerequisites

- [NPM](https://docs.npmjs.com/downloading-and-installing-node-js-and-npm)
- [Go](https://go.dev/doc/install#testing)
- [Docker](https://docs.docker.com/get-docker/)
- [Swagger](https://github.com/swaggo/swag)

## Build Go Binary
RUN ```make swagger``` to regenerate swagger docs
RUN ```make build-api``` to locally build an executable

## Tests
RUN ```make ui-test``` to test ui
RUN ```make api-test``` to test api

## Bundle react ui, api and postgres db with Docker
RUN ```make start``` to build and spin-up ui, api and db containers
RUN ```make logs``` to view logs
RUN ```make stop``` to shut down containers

- Api can be accessed at `http://localhost:8080/api/v1`
- Monitor route is at `http://localhost:8080/api/v1/monitor`
- Swagger docs can be accessed at `http://localhost:8080/swagger/index.html`

## API Routes
-  [POST] /api/v1/sensors 
-  [GET]  /api/v1/sensors/:name
-  [PUT] /api/v1/sensors/:name

## TODOs
- Better description in swagger documentation
- Unit Tests for db interfaces
- Better code comments and logs
- TLS for production build