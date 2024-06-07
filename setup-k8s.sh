#!/bin/bash

# Label nodes to identify them for specific services
# kubectl label nodes node-1 type=user-node
# kubectl label nodes node-2 type=search-node
# kubectl label nodes node-3 type=reservation-node
# kubectl label nodes node-4 type=rate-node
# kubectl label nodes node-5 type=profile-node
# kubectl label nodes node-6 type=frontend-node

# echo "Nodes have been labeled successfully."

# Apply Kubernetes YAML files for services and Redis
kubectl apply -f k8s/frontend-deployment.yaml
kubectl apply -f k8s/frontend-service.yaml

kubectl apply -f k8s/user-deployment.yaml
kubectl apply -f k8s/user-service.yaml
kubectl apply -f k8s/user-redis-deployment.yaml
kubectl apply -f k8s/user-redis-service.yaml

kubectl apply -f k8s/search-deployment.yaml
kubectl apply -f k8s/search-service.yaml
kubectl apply -f k8s/search-redis-deployment.yaml
kubectl apply -f k8s/search-redis-service.yaml

kubectl apply -f k8s/reservation-deployment.yaml
kubectl apply -f k8s/reservation-service.yaml
kubectl apply -f k8s/reservation-redis-deployment.yaml
kubectl apply -f k8s/reservation-redis-service.yaml

kubectl apply -f k8s/rate-deployment.yaml
kubectl apply -f k8s/rate-service.yaml
kubectl apply -f k8s/rate-redis-deployment.yaml
kubectl apply -f k8s/rate-redis-service.yaml

kubectl apply -f k8s/profile-deployment.yaml
kubectl apply -f k8s/profile-service.yaml
kubectl apply -f k8s/profile-redis-deployment.yaml
kubectl apply -f k8s/profile-redis-service.yaml

echo "Kubernetes resources have been applied successfully."
