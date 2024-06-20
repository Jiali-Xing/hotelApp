#!/bin/bash

# Function to find pod by label
find_pod() {
  local label=$1
  kubectl get pods -l "app=${label}" -o jsonpath="{.items[0].metadata.name}"
}

# Find the social app pods
compose_post_pod=$(find_pod "composepost")
home_timeline_pod=$(find_pod "hometimeline")
user_timeline_pod=$(find_pod "usertimeline")
social_graph_pod=$(find_pod "socialgraph")
post_storage_pod=$(find_pod "poststorage")

# Check if pods were found
if [ -z "$compose_post_pod" ]; then
  echo "Error: Compose Post pod not found"
  exit 1
fi

if [ -z "$home_timeline_pod" ]; then
  echo "Error: Home Timeline pod not found"
  exit 1
fi

if [ -z "$user_timeline_pod" ]; then
  echo "Error: User Timeline pod not found"
  exit 1
fi

if [ -z "$social_graph_pod" ]; then
  echo "Error: Social Graph pod not found"
  exit 1
fi

if [ -z "$post_storage_pod" ]; then
  echo "Error: Post Storage pod not found"
  exit 1
fi

# Establish port forwarding
echo "Forwarding localhost:50062 to ${compose_post_pod}:50051"
kubectl port-forward "${compose_post_pod}" 50062:50051 &

echo "Forwarding localhost:50059 to ${home_timeline_pod}:50051"
kubectl port-forward "${home_timeline_pod}" 50059:50051 &

echo "Forwarding localhost:50058 to ${user_timeline_pod}:50051"
kubectl port-forward "${user_timeline_pod}" 50058:50051 &

echo "Forwarding localhost:50061 to ${social_graph_pod}:50051"
kubectl port-forward "${social_graph_pod}" 50061:50051 &

# echo "Forwarding localhost:50060 to ${post_storage_pod}:50051"
# kubectl port-forward "${post_storage_pod}" 50060:50051 &

# Wait to keep script running to maintain port forwards
wait
