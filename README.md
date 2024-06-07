# hotelApp
Here we implement Hotel microservice app on grpc and redis. 

# Node Labeling Instructions

To ensure that certain pods are scheduled on specific nodes, label the nodes accordingly. Run the following commands:

```sh
# Label node-1 for the user service
kubectl label nodes node-1 type=user-node

# Label node-2 for the search service
kubectl label nodes node-2 type=search-node

# Label node-3 for the reservation service
kubectl label nodes node-3 type=reservation-node

# Label node-4 for the rate service
kubectl label nodes node-4 type=rate-node

# Label node-5 for the profile service
kubectl label nodes node-5 type=profile-node

# Label node-6 as a general-purpose node (e.g., frontend)
kubectl label nodes node-6 type=frontend-node

### Summary

1. **Run Commands**: Use the provided `kubectl` commands to label your nodes.
2. **Document**: Include the commands in your repository's documentation to ensure others know how to label the nodes.

This setup will ensure that the pods for each service are scheduled on the appropriate nodes based on the labels.
