#!/bin/bash

# Path to the msgraph.yaml file
FILE_PATH="$HOME/service-app/services/protobuf-grpc/msgraph.yaml"

# Check if the file exists
if [[ -f "$FILE_PATH" ]]; then
    # Create or update the ConfigMap
    kubectl create configmap msgraph-config --from-file=msgraph.yaml="$FILE_PATH" --dry-run=client -o yaml | kubectl apply -f -
    echo "ConfigMap msgraph-config created/updated successfully."
else
    echo "File $FILE_PATH does not exist."
    exit 1
fi
