#!/bin/bash

# inheren the environment variables DEBUG_INFO
export DEBUG_INFO
export METHOD
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

# for port in {50051..50059}; do
#   pid=$(lsof -t -i :$port)
#   if [ -n "$pid" ]; then
#     kill -9 $pid
#   fi
# done

# Exit the script
exit 0
