#!/bin/bash

echo "Testing Go module and build..."

# Test go mod
echo "Running go mod tidy..."
go mod tidy

# Test syntax
echo "Checking syntax..."
go vet ./...

echo "Build test completed!"