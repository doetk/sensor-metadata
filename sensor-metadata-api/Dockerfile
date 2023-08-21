# Use an official Golang runtime as a parent image
FROM golang:latest

# Set the working directory to /app
WORKDIR /app

# Copy the Go modules and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY .. .

# Build the application
RUN go build -o sensor-metadata-api

# Expose port 3000 to the outside world
EXPOSE 3000

# Command to run the application
CMD ["./sensor-metadata-api"]
