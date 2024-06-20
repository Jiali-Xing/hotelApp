#!/bin/bash

# Function to get the service Cluster IP
get_service_ip() {
  local service_name=$1
  kubectl get svc $service_name -o jsonpath='{.spec.clusterIP}'
}

# Get the service IPs
compose_post_ip=$(get_service_ip composepost)
home_timeline_ip=$(get_service_ip hometimeline)
user_timeline_ip=$(get_service_ip usertimeline)
social_graph_ip=$(get_service_ip socialgraph)

# Print the service IPs
echo "Compose Post IP: $compose_post_ip"
echo "Home Timeline IP: $home_timeline_ip"
echo "User Timeline IP: $user_timeline_ip"
echo "Social Graph IP: $social_graph_ip"

# Run the populate command with the correct addresses
./social-populate/populate \
  -compose_post="${compose_post_ip}:50051" \
  -home_timeline="${home_timeline_ip}:50051" \
  -user_timeline="${user_timeline_ip}:50051" \
  -social_graph="${social_graph_ip}:50051" \
  -num_of_users 100 \
  -num_of_posts 10 \
  -num_of_followers 10
