# Use the official Golang image as a build stage
FROM golang:latest as builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build all the service binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/frontend ./cmd/hotel/frontend
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/user ./cmd/hotel/user
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/search ./cmd/hotel/search
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/reservation ./cmd/hotel/reservation
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/rate ./cmd/hotel/rate
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/profile ./cmd/hotel/profile

# Use a minimal base image to run the binaries
FROM alpine:latest

# Copy the pre-built binaries from the builder stage
COPY --from=builder /app/bin /bin/

# Default to the frontend service; this can be overridden by the deployment YAML
ENTRYPOINT ["/bin/frontend"]


# # Use the official Golang image as a build stage
# FROM golang:latest

# # Set the Current Working Directory inside the container
# WORKDIR /app

# # Copy the go.mod and go.sum files
# COPY go.mod go.sum ./

# # Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download

# # Copy the source from the current directory to the Working Directory inside the container
# COPY . .

# # Build the Go app
# # RUN go build -o redis_app .


# # Command to run the executable
# CMD ["sleep", "infinity"]





# # Use the official Golang image as a build stage
# FROM golang:1.18 as builder

# # Set the Current Working Directory inside the container
# WORKDIR /app

# # Copy the go.mod and go.sum files
# COPY go.mod go.sum ./

# # Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
# RUN go mod download

# # Copy the source from the current directory to the Working Directory inside the container
# COPY . .

# # Build the Go app
# RUN go build -o redis_app .

# # Start a new stage from scratch
# FROM ubuntu:latest

# # Copy the Pre-built binary file from the previous stage
# COPY --from=builder /app/redis_app /redis_app

# # Command to run the executable
# CMD ["/redis_app"]

