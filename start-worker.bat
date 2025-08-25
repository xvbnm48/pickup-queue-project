@echo off
echo Starting Package Expiry Worker...
cd backend
go run cmd/worker/main.go
