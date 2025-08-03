#!/bin/bash

# Get the machine architecture using uname -m
arch=$(uname -m)

# Check if the architecture is 64-bit x86 or ARM
if [[ "$arch" == "x86_64" ]]; then
  echo "This is a 64-bit x86 architecture machine"
elif [[ "$arch" == "arm64" ]]; then
  echo "This is a 64-bit ARM architecture machine"
else
  echo "This is not a 64-bit system. Architecture: $arch"
fi

export CPP_IMAGE="cpp_${arch}:latest"
export GOLANG_IMAGE="golang_${arch}:latest"

# Build the server
go build -o server ../cmd/server.go

# Start the server
./server