#!/bin/bash

# Create the bin directory
mkdir -p bin

# Build each service
CGO_ENABLED=0 GOOS=linux go build -o bin/frontend ./cmd/hotel/frontend
CGO_ENABLED=0 GOOS=linux go build -o bin/user ./cmd/hotel/user
CGO_ENABLED=0 GOOS=linux go build -o bin/search ./cmd/hotel/search
CGO_ENABLED=0 GOOS=linux go build -o bin/reservation ./cmd/hotel/reservation
CGO_ENABLED=0 GOOS=linux go build -o bin/rate ./cmd/hotel/rate
CGO_ENABLED=0 GOOS=linux go build -o bin/profile ./cmd/hotel/profile
