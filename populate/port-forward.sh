#!/bin/bash

# Function to find pod by label
find_pod() {
  local label=$1
  kubectl get pods -l "app=${label}" -o jsonpath="{.items[0].metadata.name}"
}

# Find the frontend and user pods
frontend_pod=$(find_pod "frontend")
user_pod=$(find_pod "user")

# Check if frontend pod was found
if [ -z "$frontend_pod" ]; then
  echo "Error: Frontend pod not found"
  exit 1
fi

# Check if user pod was found
if [ -z "$user_pod" ]; then
  echo "Error: User pod not found"
  exit 1
fi

# Establish port forwarding
echo "Forwarding localhost:50052 to ${frontend_pod}:50051"
kubectl port-forward "${frontend_pod}" 50052:50051 &

echo "Forwarding localhost:50053 to ${user_pod}:50051"
kubectl port-forward "${user_pod}" 50053:50051 &

# Wait to keep script running to maintain port forwards
wait
