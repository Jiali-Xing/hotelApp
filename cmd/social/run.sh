#!/bin/bash

# Run Compose Post service
go run ./cmd/social/compose_post/main.go -local &

# Run Home Timeline service
go run ./cmd/social/home_timeline/main.go -local &

# Run User Timeline service
go run ./cmd/social/user_timeline/main.go -local &

# Run Social Graph service
go run ./cmd/social/social_graph/main.go -local &

# Run Post Storage service
go run ./cmd/social/post_storage/main.go -local &

# Wait for all background processes to finish
wait
