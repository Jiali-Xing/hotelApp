#!/bin/bash
kubectl delete deployment --all > /dev/null && kubectl delete service --all > /dev/null 

python ./scripts/gen-yaml.py

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

# run populate/port-forward.sh and populate/populate 
# to populate the database and establish port forwarding
./populate/port-forward.sh
./populate/populate -hotels_file=/users/jiali/hotelApp/experiments/hotel/data/hotels.json