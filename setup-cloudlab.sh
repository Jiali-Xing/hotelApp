#!/bin/bash

GHZ_VERSION="v0.120.0"
GHZ_FILENAME="ghz-linux-x86_64.tar.gz"
GHZ_URL="https://github.com/bojand/ghz/releases/download/${GHZ_VERSION}/${GHZ_FILENAME}"

# Connect to the remote node
ssh -o StrictHostKeyChecking=no -p 22 -i ${private_key} ${node_username}@${node_address} << 'EOF'

# Clone the hotelApp repository

# Function to check if a directory exists
check_dir_exists() {
    if [ -d "$1" ]; then
        echo "$1 already exists, skipping clone."
        return 1
    else
        return 0
    fi
}

# Function to download and untar ghz
download_and_untar_ghz() {
    if [ ! -f "${GHZ_FILENAME}" ]; then
        wget ${GHZ_URL}
        tar -xzf ${GHZ_FILENAME}
    else
        echo "${GHZ_FILENAME} already exists, skipping download."
    fi
}

# Add GitHub to known hosts
ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts


# Function to clone repository with SSH key verification handled
clone_repo() {
    local repo=$1
    if check_dir_exists ${repo}; then
        echo "Cloning ${repo} repository..."
        GIT_SSH_COMMAND="ssh -o StrictHostKeyChecking=accept-new" git clone git@github.com:Jiali-Xing/${repo}.git
    else
        echo "${repo} already exists."
    fi
}

# Clone the hotelApp repository
clone_repo "hotelApp"

# Clone the hotelproto repository
clone_repo "hotelproto"


# Download and untar ghz
download_and_untar_ghz

# Navigate to the hotelApp directory
cd hotelApp

# Delete all Kubernetes services, deployments, and configmaps
kubectl delete svc --all --ignore-not-found
kubectl delete deployments --all --ignore-not-found
kubectl delete configmaps --all --ignore-not-found

# Run the setup-k8s.sh script
if [ -f "./setup-k8s.sh" ]; then
    ./setup-k8s.sh
else
    echo "setup-k8s.sh not found."
fi

EOF
