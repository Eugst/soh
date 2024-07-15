# Use the official Golang image as the base image
FROM golang:1.16-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go module files
COPY go.mod ./

# Download and install the application dependencies
RUN go mod download

# Copy the source code into the container
COPY go.* main.go templates ./
COPY static ./static

# Build the Go application
RUN go build -o app

# Expose the port that the application listens on
EXPOSE 8080

# Set the entry point for the container
CMD ["./app"]
