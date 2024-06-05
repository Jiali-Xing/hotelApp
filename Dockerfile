# Use the official Golang image as a build stage
FROM golang:latest

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
# RUN go build -o redis_app .


# Command to run the executable
CMD ["sleep", "infinity"]





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

