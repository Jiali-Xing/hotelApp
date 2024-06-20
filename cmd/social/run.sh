#!/bin/bash
# Stop and remove any existing Redis containers
docker stop redis-composepost redis-hometimeline redis-usertimeline redis-socialgraph redis-poststorage
docker rm redis-composepost redis-hometimeline redis-usertimeline redis-socialgraph redis-poststorage

# Run new Redis containers
docker run -d --name redis-composepost -p 6380:6379 redis:latest
docker run -d --name redis-hometimeline -p 6381:6379 redis:latest
docker run -d --name redis-usertimeline -p 6382:6379 redis:latest
docker run -d --name redis-socialgraph -p 6383:6379 redis:latest
docker run -d --name redis-poststorage -p 6384:6379 redis:latest

# Run Compose Post service
SERVICE_NAME="composepost" REDIS_ADDR="localhost:6380" go run ./cmd/social/compose_post/main.go -local -debug &

# Run Home Timeline service
SERVICE_NAME="hometimeline" REDIS_ADDR="localhost:6381" go run ./cmd/social/home_timeline/main.go -local &

# Run User Timeline service
SERVICE_NAME="usertimeline" REDIS_ADDR="localhost:6382" go run ./cmd/social/user_timeline/main.go -local &

# Run Social Graph service
SERVICE_NAME="socialgraph" REDIS_ADDR="localhost:6383" go run ./cmd/social/social_graph/main.go -local -debug &

# Run Post Storage service
SERVICE_NAME="poststorage" REDIS_ADDR="localhost:6384" go run ./cmd/social/post_storage/main.go -local &

# Wait for all background processes to finish
wait

# # Run Compose Post service
# SERVICE_NAME="composepost" go run ./cmd/social/compose_post/main.go -local -debug &

# # Run Home Timeline service
# SERVICE_NAME="hometimeline" go run ./cmd/social/home_timeline/main.go -local &

# # Run User Timeline service
# SERVICE_NAME="usertimeline" go run ./cmd/social/user_timeline/main.go -local &

# # Run Social Graph service
# SERVICE_NAME="socialgraph" go run ./cmd/social/social_graph/main.go -local -debug &

# # Run Post Storage service
# SERVICE_NAME="poststorage" go run ./cmd/social/post_storage/main.go -local &

# # Wait for all background processes to finish
# wait
