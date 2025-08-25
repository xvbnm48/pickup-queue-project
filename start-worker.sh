#!/bin/bash

# Start Package Expiry Worker
echo "Starting Package Expiry Worker..."
cd backend
go run cmd/worker/main.go
