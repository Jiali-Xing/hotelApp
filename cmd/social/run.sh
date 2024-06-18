#!/bin/bash

# Run Compose Post service
SERVICE_NAME="composepost" go run ./cmd/social/compose_post/main.go -local -debug &

# Run Home Timeline service
SERVICE_NAME="hometimeline" go run ./cmd/social/home_timeline/main.go -local &

# Run User Timeline service
SERVICE_NAME="usertimeline" go run ./cmd/social/user_timeline/main.go -local &

# Run Social Graph service
SERVICE_NAME="socialgraph" go run ./cmd/social/social_graph/main.go -local -debug &

# Run Post Storage service
SERVICE_NAME="poststorage" go run ./cmd/social/post_storage/main.go -local &

# Wait for all background processes to finish
wait
