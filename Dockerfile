# syntax = docker/dockerfile:1.4

# Use the official Golang image as a build stage
FROM golang:latest AS builder

# Set the Current Working Directory inside the container
WORKDIR /app

# Tell Go that github.com/pennsail/* is private
ENV GOPRIVATE=github.com/pennsail/*

# Copy the go.mod and go.sum files
COPY go.mod go.sum ./

# Install git (you may already have it) and set up the https→ssh rewrite
RUN apt-get update && apt-get install -y git \
 && git config --global url."git@github.com:".insteadOf "https://github.com/"

# 3) Make sure SSH knows GitHub’s host key
RUN mkdir -p /root/.ssh \
  && ssh-keyscan github.com >> /root/.ssh/known_hosts

# Use SSH mount for private modules
RUN --mount=type=ssh \
    go mod download

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

# Build all the social network service binaries
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/composepost ./cmd/social/compose_post
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/hometimeline ./cmd/social/home_timeline
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/usertimeline ./cmd/social/user_timeline
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/socialgraph ./cmd/social/social_graph
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/poststorage ./cmd/social/post_storage
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/nginx ./cmd/social/nginx

# Use a minimal base image to run the binaries
FROM alpine:latest

# Copy the pre-built binaries from the builder stage
COPY --from=builder /app/bin /bin/
# Copy the Kubernetes config file into the image so services can load it
# bring in the YAML
COPY msgraph.yaml ./msgraph.yaml

# Default to local mode for every service
ENTRYPOINT ["/bin/sh","-c","/bin/${SERVICE_NAME} -local"]

