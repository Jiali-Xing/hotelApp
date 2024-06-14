#!/bin/bash

# inheren the environment variables DEBUG_INFO
export DEBUG_INFO
python ./scripts/gen-yaml.py

# Get all deployments and services excluding those with 'redis' in the name
deployments=$(kubectl get deployments -o custom-columns=NAME:.metadata.name --no-headers | grep -v 'redis')
services=$(kubectl get services -o custom-columns=NAME:.metadata.name --no-headers | grep -v 'redis')

sleep 5

# Delete the filtered deployments and services
for deployment in $deployments; do
  kubectl delete deployment $deployment
done

for service in $services; do
  kubectl delete service $service
done

sleep 5

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

# wait for the pods to be ready
kubectl wait --for=condition=ready pod --all --timeout=60s

# for port in {50051..50059}; do
#   pid=$(lsof -t -i :$port)
#   if [ -n "$pid" ]; then
#     kill -9 $pid
#   fi
# done

# Exit the script
exit 0
