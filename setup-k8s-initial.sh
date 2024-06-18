#!/bin/bash
kubectl delete deployment --all > /dev/null && kubectl delete service --all > /dev/null && kubectl delete configmap --all > /dev/null

sleep 5

kubectl create configmap msgraph-config --from-file=/users/jiali/hotelApp/msgraph.yaml

sleep 5
# if method is not provided, default to compose
METHOD=${1:-"compose"}
export METHOD
python scripts/gen-yaml.py

# Apply Kubernetes YAML files for services and Redis

if [[ "$METHOD" == *"hotel"* ]]; then
  # Apply Kubernetes YAML files for hotel services and Redis
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

elif [ "$METHOD" = "compose" -o "$METHOD" = "home-timeline" -o "$METHOD" = "user-timeline" -o "$METHOD" = "all-methods-social" ]; then
  # Apply Kubernetes YAML files for social network services and Redis
  kubectl apply -f k8s/composepost-deployment.yaml
  kubectl apply -f k8s/composepost-service.yaml

  kubectl apply -f k8s/hometimeline-deployment.yaml
  kubectl apply -f k8s/hometimeline-service.yaml
  kubectl apply -f k8s/hometimeline-redis-deployment.yaml
  kubectl apply -f k8s/hometimeline-redis-service.yaml

  kubectl apply -f k8s/usertimeline-deployment.yaml
  kubectl apply -f k8s/usertimeline-service.yaml
  kubectl apply -f k8s/usertimeline-redis-deployment.yaml
  kubectl apply -f k8s/usertimeline-redis-service.yaml

  kubectl apply -f k8s/socialgraph-deployment.yaml
  kubectl apply -f k8s/socialgraph-service.yaml
  kubectl apply -f k8s/socialgraph-redis-deployment.yaml
  kubectl apply -f k8s/socialgraph-redis-service.yaml

  kubectl apply -f k8s/poststorage-deployment.yaml
  kubectl apply -f k8s/poststorage-service.yaml
  kubectl apply -f k8s/poststorage-redis-deployment.yaml
  kubectl apply -f k8s/poststorage-redis-service.yaml
fi

echo "Kubernetes resources have been applied successfully."

# wait for the pods to be ready
kubectl wait --for=condition=ready pod --all --timeout=60s

for port in {50051..50059}; do
  pid=$(lsof -t -i :$port)
  if [ -n "$pid" ]; then
    kill -9 $pid
  fi
done

if [[ $METHOD == *"hotel"* ]]; then
  # Run hotel populate scripts
  ./populate/port-forward.sh &
  PORT_FORWARD_PID=$!

  # Wait for port-forwarding to be ready (adjust sleep time as needed)
  sleep 10

  ./populate/populate -hotels_file=/users/jiali/hotelApp/experiments/hotel/data/hotels.json

elif [ "$METHOD" = "compose" -o "$METHOD" = "home-timeline" -o "$METHOD" = "user-timeline" -o "$METHOD" = "all-methods-social" ]; then
  # Run social populate scripts
  ./social-populate/social-port-forward.sh &
  PORT_FORWARD_PID=$!

  # Wait for port-forwarding to be ready (adjust sleep time as needed)
  sleep 10

  ./social-populate/populate -compose_post=localhost:50062 -home_timeline=localhost:50059 -user_timeline=localhost:50058 -social_graph=localhost:50061
fi

for port in {50051..50059}; do
  pid=$(lsof -t -i :$port)
  if [ -n "$pid" ]; then
    kill -9 $pid
  fi
done

# Exit the script
exit 0
