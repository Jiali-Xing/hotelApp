#!/bin/bash
kubectl delete deployment --all > /dev/null && kubectl delete service --all > /dev/null 

python scripts/gen-yaml.py

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

for port in {50051..50059}; do
  pid=$(lsof -t -i :$port)
  if [ -n "$pid" ]; then
    kill -9 $pid
  fi
done

# run populate/port-forward.sh and populate/populate 
# to populate the database and establish port forwarding
./populate/port-forward.sh &
PORT_FORWARD_PID=$!

# Wait for port-forwarding to be ready (adjust sleep time as needed)
sleep 10

./populate/populate -hotels_file=/users/jiali/hotelApp/experiments/hotel/data/hotels.json

# After populate script finishes, kill the port-forwarding process

for port in {50051..50059}; do
  pid=$(lsof -t -i :$port)
  if [ -n "$pid" ]; then
    kill -9 $pid
  fi
done

# Exit the script
exit 0
