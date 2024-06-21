# Start from a base image containing the Go runtime
FROM golang:1.22.4-alpine AS builder

# Set the current working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code from the current directory to the working directory inside the container
COPY . .

# Build the Go app
RUN go build -o myapp

# Start a new stage from scratch
FROM alpine:latest

# Set the working directory in the container
WORKDIR /app

# Copy the compiled binary from the builder stage to the new stage
COPY --from=builder /app/myapp .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
CMD ["./myapp"]
