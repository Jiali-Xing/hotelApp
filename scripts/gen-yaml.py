import os

printScreen = False

# Define the services for the hotel app
hotel_services = ["frontend", "user", "search", "reservation", "rate", "profile"]
hotel_services_using_redis = ["user", "search", "reservation", "rate", "profile"]

# Define the services for the social network app
social_services = ["composepost", "hometimeline", "usertimeline", "socialgraph", "poststorage", "nginx"]
social_services_using_redis = ["hometimeline", "usertimeline", "poststorage", "socialgraph"]  # Assuming these services use Redis for state management

node_list = ["node-1", "node-2", "node-3", "node-4", "node-5", "node-6"]

# if  [ "$METHOD" = "compose" -o "$METHOD" = "home-timeline" -o "$METHOD" = "user-timeline" -o "$METHOD" = "all-methods-social" ]; then
method = os.getenv("METHOD", "social")
services = hotel_services if 'hotel' in method else social_services
services_using_redis = hotel_services_using_redis if 'hotel' in method else social_services_using_redis

# Define the base directory for the output
output_dir = "k8s"

# Load the templates
with open("scripts/deploy_template.yaml", "r") as template_file:
    deploy_template = template_file.read()
with open("scripts/service_template.yaml", "r") as template_file:
    service_template = template_file.read()

with open("scripts/redis_service_template.yaml", "r") as template_file:
    redis_service_template = template_file.read()
with open("scripts/redis_deployment_template.yaml", "r") as template_file:
    redis_deployment_template = template_file.read()

# Check for DEBUG_INFO environment variable
debug_info = os.getenv("DEBUG_INFO", "false").lower() == "true"

# Generate the deployment and service YAML for each service
for index, service in enumerate(services):
    # Assign node in order, wrapping around if necessary
    node_index = index % len(node_list)
    node_name = node_list[node_index]

    if printScreen:
        args = 'args: ["/bin/{} -debug"]'.format(service)
    elif debug_info:
        args = 'args: ["/bin/{} -debug 2>&1 | tee /root/deathstar_{}.output"]'.format(service, service)
    else:
        args = 'args: ["/bin/{} 2>&1 | tee /root/deathstar_{}.output"]'.format(service, service)
    
    deployment_content = deploy_template.format(service_name=service, args=args, node_name=node_name)
    deployment_filename = "{}-deployment.yaml".format(service)
    
    with open(os.path.join(output_dir, deployment_filename), "w") as f:
        f.write(deployment_content)

    external_ip = ""

    # Conditionally add an additional port for nginx or frontend
    if service in ["nginx", "frontend"]:
        additional_port = """
  - name: http-alt  # Name for the additional port
    protocol: TCP
    port: 8082
    targetPort: 8082"""
    else:
        additional_port = ""
    
    service_content = service_template.format(service_name=service, external_ip=external_ip, additional_port=additional_port)
    service_filename = "{}-service.yaml".format(service)
    with open(os.path.join(output_dir, service_filename), "w") as f:
        f.write(service_content)

    # Generate the redis service and deployment YAML if the service uses Redis
    if service in services_using_redis:
        redis_service_content = redis_service_template.format(service_name=service)
        redis_service_filename = "{}-redis-service.yaml".format(service)
        with open(os.path.join(output_dir, redis_service_filename), "w") as f:
            f.write(redis_service_content)
        
        redis_deployment_content = redis_deployment_template.format(service_name=service, node_name=node_name)
        redis_deployment_filename = "{}-redis-deployment.yaml".format(service)
        with open(os.path.join(output_dir, redis_deployment_filename), "w") as f:
            f.write(redis_deployment_content)

print("Kubernetes YAML files have been generated in the 'k8s' directory.")
