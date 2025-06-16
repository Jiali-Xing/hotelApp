# Single‐Service Profiling Guide

This README shows how to spin up a single backend service (with colocated Redis) and drive it with **ghz** for latency/throughput characterization. We cover:

1. **PostStorage** service  
2. **SocialGraph** service  
3. Cleanup of running services  

---

## Prerequisites

- Go toolchain (to run the service)  
- Docker (for Redis)  
- Other required files are in this folder.

## 1. PostStorage Service

### 1.1 Start Redis  
```bash
docker run -d --name redis-poststorage \
  -p 6384:6379 \
  redis:latest
```

### 1.2 Run the Go server

From the **hotelApp** repo root:

```bash
SERVICE_NAME="poststorage" \
REDIS_ADDR="localhost:6384" \
go run ./cmd/social/post_storage/main.go -local &
```

> Listens on port `50060` by default.

### 1.3 Store posts (batch of 10)

In your **ghz** folder (`cmd/ghz`):

```bash
./ghz \
  --insecure \
  --proto post_storage.proto \
  --call socialproto.PostStorage.StorePostMulti \
  --data '{"creator_id":"{{.RequestNumber}}","text":"{{.RequestNumber}}","number":10}' \
  --concurrency 100 \
  --total 99999 \
  --output results/StorePostMulti.json \
  localhost:50060
```

* `{{.RequestNumber}}` → deterministic ID/text per request
* Adjust `--concurrency` / `--total` as needed
* Results (latencies, RPS) in `results/StorePostMulti.json`

### 1.4 Read back posts

```bash
./ghz \
  --insecure \
  --proto post_storage.proto \
  --call socialproto.PostStorage.ReadPosts \
  --data '{"post_ids":["{{.RequestNumber}}_{{.RequestNumber}}"]}' \
  --concurrency 100 \
  --total 99999 \
  --output results/ReadPosts.json \
  localhost:50060
```

* Reads the same deterministic IDs generated above
* Sweeps RPS across 100 concurrent streams

---

## 2. SocialGraph Service

### 2.1 Start Redis

```bash
docker run -d --name redis-socialgraph \
  -p 6390:6379 \
  redis:latest
```

### 2.2 Run the Go server

```bash
SERVICE_NAME="socialgraph" \
REDIS_ADDR="localhost:6390" \
go run ./cmd/social/social_graph/main.go -local &
```

> Listens on port `50061` by default.

### 2.3 Insert users

```bash
./ghz \
  --insecure \
  --proto social_graph.proto \
  --call socialproto.SocialGraph.InsertUser \
  --concurrency 10 \
  --total 100 \
  --data '{"user_id":"user{{.RequestNumber}}"}' \
  localhost:50061
```

Creates `user1` through `user100`.

### 2.4 Build follow relationships

```bash
./ghz \
  --insecure \
  --proto social_graph.proto \
  --call socialproto.SocialGraph.Follow \
  --concurrency 10 \
  --total 99 \
  --data '{
    "follower_id":"user{{.RequestNumber}}",
    "followee_id":"user{{add .RequestNumber 1}}"
  }' \
  localhost:50061
```

Makes `user1 → user2`, `user2 → user3`, …, `user99 → user100`.

### 2.5 Query followers

```bash
./ghz \
  --insecure \
  --proto social_graph.proto \
  --call socialproto.SocialGraph.GetFollowers \
  --concurrency 10 \
  --total 1000 \
  --data '{
    "user_id":"user{{ add (mod (sub .RequestNumber 1) 99) 1 }}"
  }' \
  localhost:50061
```

Cycles through `user1`…`user99` across 1 000 requests.
*(Replace `GetFollowers` with `GetFollowees` to query followees.)*

---

## 3. Cleanup

To kill all Go servers listening on ports `50051`–`50069`:

```bash
for port in {50051..50069}; do
  pid=$(lsof -t -i :$port) || continue
  kill -9 $pid
done
docker rm -f redis-poststorage redis-socialgraph
```

---

**Now you have a repeatable recipe** to:

* Stand up a single backend service + Redis
* Hammer it with **ghz** using deterministic templating
* Sweep concurrency, batch sizes, and total RPS
* Collect and store per‐request latency & throughput data for downstream modeling or power‐profiling.

---
## Profile the Backend Service in Docker

To profile the backend service in Docker, you can use the following steps:

0. Create a Docker network:
```bash
docker network create char-net
```
1. Setup the redis container:
```bash
docker run -d \                                
  --network char-net \     
  --name redis-poststorage \
  -p 6384:6379 \
  redis:latest
```
2. Setup the poststorage service in Docker:
```bash
docker run -d \
  --network char-net \
  --name poststorage \
  -v "$(pwd)/msgraph.yaml:/app/msgraph.yaml:ro" --entrypoint /bin/poststorage \
  -e REDIS_ADDR="redis-poststorage:6379" \
  -e MSGRAPH_YAML="/app/msgraph.yaml" \
  -p 50060:50060 \
  xjiali/social-hotel:char
```
3. Run ghz to profile the service as previously:
```bash
 # Store a batch of 10 posts                    
./ghz \
  --insecure \
  --proto post_storage.proto \
  --call socialproto.PostStorage.StorePostMulti \
  --data '{"creator_id":"{{.RequestNumber}}","text":"{{.RequestNumber}}","number":10}' \
  --concurrency 100 \
  --total 99999 \
  localhost:50060 \
  --output results/StorePostMulti.json

# Read them back
./ghz \
  --insecure \
  --proto post_storage.proto \
  --call socialproto.PostStorage.ReadPosts \
  --data '{"post_ids":["{{.RequestNumber}}_{{.RequestNumber}}"]}' \
  --concurrency 100 \
  --total 99999 \
  localhost:50060 \
  --output results/ReadPosts.json
```
4. Cleanup:
```bash
docker rm -f poststorage redis-poststorage
docker network rm char-net
```
This will allow you to profile the backend service running in Docker, ensuring that it can be easily set up and torn down without affecting your local environment.