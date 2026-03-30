# !/bin/bash
# This script is used to setup the k8s cluster for the hotel app, either for hotel or social network
# Jiali: 1. delete all deployments, services, and configmaps
# Jiali: 2. create a new configmap msgraph-config that contains no interceptors
# Jiali: 3. generate the yaml files for the services and redis
# Jiali: 4. apply the yaml files for the services and redis
# Jiali: 5. run the populate scripts for the hotel app or social network
# Jiali: 6. delete all deployments, services, and configmaps

# example: run `./setup-k8s-initial.sh` to setup the k8s cluster for the social network
# default method is compose
# example: run `./setup-k8s-initial.sh hotel` to setup the k8s cluster for the hotel app

kubectl delete deployment --all > /dev/null && kubectl delete service --all > /dev/null && kubectl delete configmap --all > /dev/null

sleep 1

# # if configmap is not provided, create a new one
# if [ ! kubectl get configmap msgraph-config ]; then
#   kubectl create configmap msgraph-config --from-file=/users/jiali/hotelApp/msgraph.yaml
# fi
echo "Creating configmap msgraph-config that contains no interceptors"
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
  kubectl apply -f k8s/nginx-deployment.yaml
  kubectl apply -f k8s/nginx-service.yaml

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

kubectl label nodes node-1 service=frontend
kubectl label nodes node-2 service=user
kubectl label nodes node-3 service=search
kubectl label nodes node-4 service=reservation
kubectl label nodes node-5 service=rate
kubectl label nodes node-6 service=profile


for port in {50051..50059}; do
  pid=$(lsof -t -i :$port)
  if [ -n "$pid" ]; then
    kill -9 $pid
  fi
done

if [[ $METHOD == *"hotel"* ]]; then
  # Run hotel populate scripts
  # ./populate/port-forward.sh &
  # PORT_FORWARD_PID=$!

  # Wait for port-forwarding to be ready (adjust sleep time as needed)
  # sleep 25

  # # Query the frontend and user addresses from user input
  # read -p "Enter the frontend address: " frontend_address
  # read -p "Enter the user address: " user_address

  # Get the names of all deployments
  deployments=$(kubectl get deployments -o custom-columns=":metadata.name" --no-headers)
  echo "Deployments: $deployments"

  # Wait for all deployments to be ready
  for deployment in $deployments; do
    kubectl rollout status deployment/$deployment
  done

  # Query the frontend and user service addresses automatically from Kubernetes
  sleep 25

  # Replace these with your actual service names
  frontend_service_name="frontend"
  user_service_name="user"

  # Get the Cluster IP and NodePort of the frontend service
  frontend_ip=$(kubectl get service $frontend_service_name -o=jsonpath='{.spec.clusterIP}')
  # frontend_nodeport=$(kubectl get service $frontend_service_name -o=jsonpath='{.spec.ports[0].nodePort}')
  frontend_address="$frontend_ip:50051"

  # Get the Cluster IP and NodePort of the user service
  user_ip=$(kubectl get service $user_service_name -o=jsonpath='{.spec.clusterIP}')
  # user_nodeport=$(kubectl get service $user_service_name -o=jsonpath='{.spec.ports[0].nodePort}')
  user_address="$user_ip:50051"

  echo "Frontend Address: $frontend_address"
  echo "User Address: $user_address"

  ./populate/populate -frontend="$frontend_address" -user="$user_address" -num_of_users=1000 -hotels_file=/users/jiali/hotelApp/experiments/hotel/data/hotels.json > popu.output 2>&1
  # -hotels_file=/users/jiali/hotelApp/experiments/hotel/data/hotels.json -num_of_users=1000 > popu.output 2>&1

elif [ "$METHOD" = "compose" -o "$METHOD" = "home-timeline" -o "$METHOD" = "user-timeline" -o "$METHOD" = "all-methods-social" ]; then
  # Run social populate scripts
  # ./social-populate/social-port-forward.sh &
  # PORT_FORWARD_PID=$!

  # Wait for port-forwarding to be ready (adjust sleep time as needed)
  sleep 55

  ./social-populate/populate.sh > popu.output 2>&1 
fi

for port in {50051..50059}; do
  pid=$(lsof -t -i :$port)
  if [ -n "$pid" ]; then
    kill -9 $pid
  fi
done

# Exit the script
exit 0
