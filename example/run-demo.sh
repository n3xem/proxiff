#!/bin/bash

# Proxiff Demo Script
# This script starts sample servers and proxiff for demonstration

set -e

echo "Building proxiff..."
go build -o proxiff ./cmd/proxiff

echo "Building sample server..."
go build -o sample-server ./example/servers

echo ""
echo "Starting demo environment..."
echo "================================"
echo ""

# Start current server
echo "Starting current server on :8081..."
./sample-server -port 8081 -version current &
CURRENT_PID=$!

# Start newer server
echo "Starting newer server on :8082..."
./sample-server -port 8082 -version newer &
NEWER_PID=$!

# Wait for servers to start
sleep 2

# Start proxiff
echo "Starting proxiff on :8080..."
./proxiff -newer http://localhost:8082 -current http://localhost:8081 -port 8080 &
PROXY_PID=$!

echo ""
echo "================================"
echo "Demo environment is ready!"
echo "================================"
echo ""
echo "Current server:  http://localhost:8081"
echo "Newer server:    http://localhost:8082"
echo "Proxiff:         http://localhost:8080"
echo ""
echo "Try these commands:"
echo "  curl http://localhost:8080/"
echo "  curl http://localhost:8080/api/users"
echo "  curl http://localhost:8080/api/status"
echo ""
echo "Press Ctrl+C to stop all servers"
echo ""

# Cleanup function
cleanup() {
    echo ""
    echo "Stopping servers..."
    kill $CURRENT_PID $NEWER_PID $PROXY_PID 2>/dev/null || true
    echo "Done!"
    exit 0
}

trap cleanup INT TERM

# Wait for user to press Ctrl+C
wait
