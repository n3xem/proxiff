#!/bin/bash

set -e

# Check if protoc is installed
if ! command -v protoc &> /dev/null; then
    echo "Error: protoc is not installed"
    echo "Please install protoc: https://grpc.io/docs/protoc-installation/"
    exit 1
fi

# Generate Go code from proto files
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    comparator.proto

echo "Proto files generated successfully!"
