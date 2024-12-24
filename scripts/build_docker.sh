#!/bin/bash

# Get the machine architecture using uname -m
arch=$(uname -m)

# Set the Docker image version to latest
version="latest"

# Check if the architecture is 64-bit x86 or ARM
if [[ "$arch" == "x86_64" ]]; then
  echo "This is a 64-bit x86 architecture machine"
elif [[ "$arch" == "arm64" ]]; then
  echo "This is a 64-bit ARM architecture machine"
else
  echo "This is not a 64-bit system. Architecture: $arch"
fi

script_dir=$(dirname "$(realpath "$0")")
dockerfiles_folder="$script_dir/../Dockerfiles/$arch/"

# Define the dockerfile_map as a regular associative array in zsh
languages=("cpp" "golang")
dockerfiles=("cpp.Dockerfile" "golang.Dockerfile")

# Loop through each Dockerfile and build the Docker image
for i in "${!languages[@]}"; do
    language=${languages[$i]}
    dockerfile="${dockerfiles_folder}${dockerfiles[$i]}"
    if [[ -f "$dockerfile" ]]; then
        echo "Building Docker image for $language using $dockerfile"
        docker build -t "${language}_${arch}:latest" -f "$dockerfile" "${dockerfiles_folder}"
    else
        echo "Dockerfile for $language not found: $dockerfile"
    fi
done
